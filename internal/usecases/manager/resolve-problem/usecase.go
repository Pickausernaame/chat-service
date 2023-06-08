package resolveproblem

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	problemsrepo "github.com/Pickausernaame/chat-service/internal/repositories/problems"
	jobresolveproblem "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/job-resolve-problem"
	"github.com/Pickausernaame/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=resolveproblemmocks

var (
	ErrInvalidRequest  = errors.New("invalid request")
	ErrProblemNotFound = errors.New("problem not found")
)

type messagesRepository interface {
	CreateProblemResolvedMessage(ctx context.Context, chatID types.ChatID, problemID types.ProblemID, reqID types.RequestID) (*messagesrepo.Message, error)
}

type problemsRepository interface {
	GetProblemByChatAndManagerIDs(ctx context.Context, chatID types.ChatID, managerID types.UserID) (*problemsrepo.Problem, error)
	ResolveProblem(ctx context.Context, problemID types.ProblemID, managerID types.UserID) error
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

type outboxService interface {
	Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	outbox  outboxService      `option:"mandatory" validate:"required"`
	prbRepo problemsRepository `option:"mandatory" validate:"required"`
	msgRepo messagesRepository `option:"mandatory" validate:"required"`
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
	return UseCase{Options: opts, lg: zap.L().Named("resolve-problem")}, nil
}

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if err := req.Validate(); err != nil {
		return Response{}, fmt.Errorf("request validation: %v %w", err, ErrInvalidRequest)
	}

	p, err := u.prbRepo.GetProblemByChatAndManagerIDs(ctx, req.ChatID, req.ManagerID)
	if err != nil {
		if errors.Is(err, problemsrepo.ErrProblemNotFound) {
			return Response{}, ErrProblemNotFound
		}
		return Response{}, fmt.Errorf("getting problem by chat id: %v", err)
	}
	err = u.txtr.RunInTx(ctx, func(ctx context.Context) error {
		msg, err := u.msgRepo.CreateProblemResolvedMessage(ctx, req.ChatID, p.ID, req.RequestID)
		if err != nil {
			return fmt.Errorf("creating problem resolved message: %v", err)
		}

		err = u.prbRepo.ResolveProblem(ctx, p.ID, req.ManagerID)
		if err != nil {
			return fmt.Errorf("resolving problem: %w", err)
		}

		payload, err := jobresolveproblem.MarshalPayload(req.ManagerID, req.RequestID, msg.ID, req.ChatID)
		if err != nil {
			return fmt.Errorf("marshaling payload: %v", err)
		}
		_, err = u.outbox.Put(ctx, jobresolveproblem.Name, payload, time.Now())
		if err != nil {
			return fmt.Errorf("put outbox job: %v", err)
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, problemsrepo.ErrProblemNotFound) {
			return Response{}, ErrProblemNotFound
		}
		return Response{}, fmt.Errorf("tx failed: %v", err)
	}

	return Response{}, nil
}
