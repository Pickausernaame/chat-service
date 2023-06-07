package sendmessage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	problemsrepo "github.com/Pickausernaame/chat-service/internal/repositories/problems"
	sendmanagermessagejob "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/send-manager-message"
	"github.com/Pickausernaame/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=sendmessagemocks

var ErrInvalidRequest = errors.New("invalid request")

type messagesRepository interface {
	CreateFullVisible(
		ctx context.Context,
		reqID types.RequestID,
		problemID types.ProblemID,
		chatID types.ChatID,
		authorID types.UserID,
		msgBody string,
	) (*messagesrepo.Message, error)
}

type outboxService interface {
	Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
}

type problemsRepository interface {
	GetProblemByChatAndManagerIDs(ctx context.Context, chatID types.ChatID, managerID types.UserID) (*problemsrepo.Problem, error)
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	msgRepo messagesRepository `option:"mandatory" validate:"required"`
	outbox  outboxService      `option:"mandatory" validate:"required"`
	prbRepo problemsRepository `option:"mandatory" validate:"required"`
	txtr    transactor         `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
	lg *zap.Logger
}

func New(opts Options) (UseCase, error) {
	if err := opts.Validate(); err != nil {
		return UseCase{}, fmt.Errorf("validating: %v", err)
	}
	return UseCase{Options: opts, lg: zap.L().Named("send-message-usecase")}, nil
}

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if err := req.Validate(); err != nil {
		return Response{}, fmt.Errorf("request validation: %v %w", err, ErrInvalidRequest)
	}

	p, err := u.prbRepo.GetProblemByChatAndManagerIDs(ctx, req.ChatID, req.ManagerID)
	if err != nil {
		return Response{}, fmt.Errorf("getting problem by chatID and managerID: %v", err)
	}

	msg := &messagesrepo.Message{}
	err = u.txtr.RunInTx(ctx, func(ctx context.Context) error {
		msg, err = u.msgRepo.CreateFullVisible(ctx, req.ID, p.ID, req.ChatID, req.ManagerID, req.MessageBody)
		if err != nil {
			return fmt.Errorf("creating full visible msg: %v", err)
		}

		payload, err := sendmanagermessagejob.MarshalPayload(msg.ID)
		if err != nil {
			return fmt.Errorf("marshaling payload: %v", err)
		}
		_, err = u.outbox.Put(ctx, sendmanagermessagejob.Name, payload, time.Now())
		if err != nil {
			return fmt.Errorf("put outbox job: %v", err)
		}
		return nil
	})
	if err != nil {
		return Response{}, fmt.Errorf("tx failed: %v", err)
	}

	return Response{
		MessageID: msg.ID,
		CreatedAt: msg.CreatedAt,
	}, nil
}
