package afcverdictsprocessor

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/golang-jwt/jwt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	clientmessageblockedjob "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/client-message-blocked"
	clientmessagesentjob "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/client-message-sent"
	"github.com/Pickausernaame/chat-service/internal/types"
	"github.com/Pickausernaame/chat-service/internal/validator"
)

const (
	serviceName           = "afc-verdict-processor"
	suspiciousVerdictType = "suspicious"
	okVerdictType         = "ok"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/service_mock.gen.go -package=afcverdictsprocessormocks

type messagesRepository interface {
	MarkAsVisibleForManager(ctx context.Context, msgID types.MessageID) error
	BlockMessage(ctx context.Context, msgID types.MessageID) error
}

type outboxService interface {
	Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

//go:generate options-gen -out-filename=service_options.gen.go -from-struct=Options
type Options struct {
	backoffInitialInterval time.Duration `default:"100ms" validate:"min=50ms,max=1s"`
	backoffMaxElapsedTime  time.Duration `default:"5s" validate:"min=500ms,max=1m"`

	brokers          []string `option:"mandatory" validate:"min=1"`
	consumers        int      `option:"mandatory" validate:"min=1,max=16"`
	consumerGroup    string   `option:"mandatory" validate:"required"`
	verdictsTopic    string   `option:"mandatory" validate:"required"`
	processBatchSize int
	verdictsSignKey  string

	readerFactory KafkaReaderFactory `option:"mandatory" validate:"required"`
	dlqWriter     KafkaDLQWriter     `option:"mandatory" validate:"required"`

	txtor   transactor         `option:"mandatory" validate:"required"`
	msgRepo messagesRepository `option:"mandatory" validate:"required"`
	outBox  outboxService      `option:"mandatory" validate:"required"`
}

type Service struct {
	lg     *zap.Logger
	pubKey *rsa.PublicKey
	Options
}

func New(opts Options) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validating opts: %v", err)
	}

	var pubKey *rsa.PublicKey
	if opts.verdictsSignKey != "" {
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(opts.verdictsSignKey))
		if err != nil {
			return nil, fmt.Errorf("parsing rsa: %v", err)
		}
		pubKey = key
	}

	if opts.processBatchSize == 0 {
		opts.processBatchSize = 1
	}

	return &Service{
		lg:      zap.L().Named(serviceName),
		Options: opts,
		pubKey:  pubKey,
	}, nil
}

func (s *Service) Run(ctx context.Context) error {
	defer func() {
		if err := s.dlqWriter.Close(); err != nil {
			s.lg.Error("closing dqlPublisher error", zap.Error(err))
		}
	}()

	eg, ctx := errgroup.WithContext(ctx)
	for i := 0; i < s.consumers; i++ {
		eg.Go(func() error {
			consumer := s.readerFactory(s.brokers, s.consumerGroup, s.verdictsTopic)
			defer func() {
				if err := consumer.Close(); err != nil {
					s.lg.Error("closing consumer error", zap.Error(err))
				}
			}()
			wg := sync.WaitGroup{}
			for {
				select {
				case <-ctx.Done():
					return nil
				default:

					var msgs []kafka.Message
					for i := 0; i < s.processBatchSize; i++ {
						ctx, cancel := context.WithTimeout(ctx, time.Second*3)
						msg, err := consumer.FetchMessage(ctx)
						cancel()
						if err != nil {
							if errors.Is(err, context.DeadlineExceeded) {
								continue
							}
							s.lg.Error("fetch message error", zap.Error(err))
							return err
						}
						msgs = append(msgs, msg)
					}

					for _, msg := range msgs {
						v, err := s.extractVerdict(msg.Value)
						if err != nil {
							s.lg.Error("extract verdict error", zap.Error(err))
							wg.Add(1)
							go func(msg kafka.Message, err error) {
								defer wg.Done()
								err = s.dlqWriter.WriteMessages(ctx, prepareDLQMessage(msg, err))
								if err != nil {
									s.lg.Error("dlqWriter error", zap.Error(err))
								}
							}(msg, err)
							continue
						}

						switch v.Status {
						case suspiciousVerdictType:
							s.lg.Error("verdict is suspicious", zap.Any("verdict", v))
							err := backoff.Retry(s.handleSuspicious(ctx, v.MessageID), backoff.WithContext(s.newBackoff(), ctx))
							if err != nil {
								s.lg.Error("handle suspicious error", zap.Error(err))
								wg.Add(1)
								go func(msg kafka.Message, err error) {
									defer wg.Done()
									err = s.dlqWriter.WriteMessages(ctx, prepareDLQMessage(msg, err))
									if err != nil {
										s.lg.Error("dlqWriter error", zap.Error(err))
									}
								}(msg, err)
								continue
							}
						case okVerdictType:
							err := backoff.Retry(s.handleOk(ctx, v.MessageID), backoff.WithContext(s.newBackoff(), ctx))
							if err != nil {
								s.lg.Error("retry with exponential backoff error", zap.Error(err))
								wg.Add(1)
								go func(msg kafka.Message, err error) {
									defer wg.Done()
									err = s.dlqWriter.WriteMessages(ctx, prepareDLQMessage(msg, err))
									if err != nil {
										s.lg.Error("dlqWriter error", zap.Error(err))
									}
								}(msg, err)
								continue
							}
						}
					}

					wg.Wait()
					err := consumer.CommitMessages(ctx, msgs...)
					if err != nil {
						s.lg.Error("commit msg error", zap.Error(err))
						return err
					}
				}
			}
		})
	}
	err := eg.Wait()
	if err != nil {
		s.lg.Error("waiting error", zap.Error(err))
	}

	return nil
}

func (s *Service) extractVerdict(msg []byte) (*verdict, error) {
	v := &verdict{}
	if s.pubKey != nil {
		cl := &jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(string(msg), cl, func(token *jwt.Token) (interface{}, error) {
			return s.pubKey, nil
		})
		if err != nil {
			return nil, fmt.Errorf("parsing token: %v", err)
		}

		msg, err = json.Marshal(cl)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal claims: %v", err)
		}
	}

	if err := json.Unmarshal(msg, v); err != nil {
		return nil, fmt.Errorf("unmarshaling payload to verdict: %v", err)
	}

	if err := validator.Validator.Struct(*v); err != nil {
		return nil, fmt.Errorf("validating verdict: %v", err)
	}

	return v, nil
}

func (s *Service) handleSuspicious(ctx context.Context, msgID types.MessageID) func() error {
	return func() error {
		return s.txtor.RunInTx(ctx, func(ctx context.Context) error {
			err := s.msgRepo.BlockMessage(ctx, msgID)
			if err != nil {
				return fmt.Errorf("blocking msg: %v", err)
			}
			payload, err := clientmessageblockedjob.MarshalPayload(msgID)
			if err != nil {
				return fmt.Errorf("marshaling clientmessageblockedjob payload: %v", err)
			}
			_, err = s.outBox.Put(ctx, clientmessageblockedjob.Name, payload, time.Now())
			if err != nil {
				return fmt.Errorf("putting %s job: %v", clientmessageblockedjob.Name, err)
			}
			return nil
		})
	}
}

func (s *Service) handleOk(ctx context.Context, msgID types.MessageID) func() error {
	return func() error {
		return s.txtor.RunInTx(ctx, func(ctx context.Context) error {
			err := s.msgRepo.MarkAsVisibleForManager(ctx, msgID)
			if err != nil {
				return fmt.Errorf("marking visible for manager msg: %v", err)
			}

			payload, err := clientmessagesentjob.MarshalPayload(msgID)
			if err != nil {
				return fmt.Errorf("marshaling clientmessagesentjob payload: %v", err)
			}

			_, err = s.outBox.Put(ctx, clientmessagesentjob.Name, payload, time.Now())
			if err != nil {
				return fmt.Errorf("putting %s job: %v", clientmessagesentjob.Name, err)
			}
			return nil
		})
	}
}

func (s *Service) newBackoff() *backoff.ExponentialBackOff {
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = s.backoffInitialInterval
	bo.MaxElapsedTime = s.backoffMaxElapsedTime
	bo.Reset()
	return bo
}
