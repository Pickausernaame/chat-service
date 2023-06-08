package managerscheduler

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	problemsrepo "github.com/Pickausernaame/chat-service/internal/repositories/problems"
	managerpool "github.com/Pickausernaame/chat-service/internal/services/manager-pool"
	managerassignedtoproblemjob "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/manager-assigned-to-problem"
	"github.com/Pickausernaame/chat-service/internal/store"
	"github.com/Pickausernaame/chat-service/internal/types"
)

const serviceName = "manager-scheduler"

type problemRepository interface {
	GetUnassignedProblems(ctx context.Context) ([]*problemsrepo.Problem, error)
	AssignManager(ctx context.Context, problemID types.ProblemID, managerID types.UserID) error
	ResolveProblem(ctx context.Context, problemID types.ProblemID, managerID types.UserID) error
}

type messageRepository interface {
	CreateProblemAssignedMessage(ctx context.Context, id types.ChatID,
		managerID types.UserID, problemID types.ProblemID) (*messagesrepo.Message, error)
	MessageForManagerByChatID(ctx context.Context, id types.ChatID) (*messagesrepo.Message, error)
}

type outbox interface {
	Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	period time.Duration `option:"mandatory" validate:"min=100ms,max=1m"`

	mngrPool managerpool.Pool `option:"mandatory" validate:"required"`

	msgRepo messageRepository `option:"mandatory" validate:"required"`

	outbox outbox `option:"mandatory" validate:"required"`

	prbRepo problemRepository `option:"mandatory" validate:"required"`

	txtr transactor `option:"mandatory" validate:"required"`
}

type Service struct {
	Options
	lg *zap.Logger
}

func New(opts Options) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validating opts: %v", err)
	}
	return &Service{
		Options: opts,
		lg:      zap.L().Named(serviceName),
	}, nil
}

func (s *Service) Run(ctx context.Context) (err error) {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:

			ps, err := s.prbRepo.GetUnassignedProblems(ctx)
			if err != nil {
				s.lg.Error("getting unassigned problems", zap.Error(err))
				continue
			}

			for _, p := range ps {
				clientMsg, err := s.msgRepo.MessageForManagerByChatID(ctx, p.ChatID)
				if err != nil {
					if store.IsNotFound(err) {
						// mark problem solved, it is zombie-problem
						if err = s.prbRepo.ResolveProblem(ctx, p.ID, types.UserIDNil); err != nil {
							s.lg.Error("resolving problem by ID", zap.Error(err))
						}
					}
					s.lg.Error("checking existing messages by chatID", zap.Error(err), zap.Stringer("chatID", p.ChatID))
					continue
				}

				managerID, err := s.mngrPool.Get(ctx)
				if err != nil {
					s.lg.Error("getting manager pool", zap.Error(err))
					continue
				}

				err = s.txtr.RunInTx(ctx, func(ctx context.Context) error {
					if err = s.prbRepo.AssignManager(ctx, p.ID, managerID); err != nil {
						return fmt.Errorf("assign manager: %v", err)
					}

					msg, err := s.msgRepo.CreateProblemAssignedMessage(ctx, p.ChatID, managerID, p.ID)
					if err != nil {
						return fmt.Errorf("creating problem assigned message: %v", err)
					}

					payload, err := managerassignedtoproblemjob.MarshalPayload(clientMsg.AuthorID, managerID, clientMsg.InitialRequestID, msg.ID)
					if err != nil {
						return fmt.Errorf("marshaling managerassignedtoproblemjob payload: %v", err)
					}

					_, err = s.outbox.Put(ctx, managerassignedtoproblemjob.Name, payload, time.Now())
					if err != nil {
						return fmt.Errorf("outbox put job: %v", err)
					}
					return nil
				})
				if err != nil {
					s.lg.Error("process problem", zap.Error(err))
					continue
				}
			}
			time.Sleep(s.period)
		}
	}
}
