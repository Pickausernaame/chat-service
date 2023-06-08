package getassignedproblems

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	problemsrepo "github.com/Pickausernaame/chat-service/internal/repositories/problems"
	"github.com/Pickausernaame/chat-service/internal/types"
)

var ErrInvalidRequest = errors.New("invalid request")

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=getassignedproblemsmocks

type problemRepository interface {
	GetAssignedUnsolvedProblems(ctx context.Context, managerID types.UserID) ([]*problemsrepo.Problem, error)
}

type chatRepository interface {
	ClientIDByID(ctx context.Context, id types.ChatID) (types.UserID, error)
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	problemRepo problemRepository `option:"mandatory" validate:"required"`
	chatRepo    chatRepository    `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
	lg *zap.Logger
}

func New(opts Options) (UseCase, error) {
	if err := opts.Validate(); err != nil {
		return UseCase{}, fmt.Errorf("validating opts: %v", err)
	}
	return UseCase{
		Options: opts,
		lg:      zap.L().Named("can-reseive-problems-usecase"),
	}, nil
}

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if err := req.Validate(); err != nil {
		return Response{}, ErrInvalidRequest
	}

	ps, err := u.problemRepo.GetAssignedUnsolvedProblems(ctx, req.ManagerID)
	if err != nil {
		return Response{}, fmt.Errorf("getting assigned problems: %v", err)
	}

	res := Response{Chats: make([]*Chat, 0, len(ps))}
	for _, p := range ps {
		clientID, err := u.chatRepo.ClientIDByID(ctx, p.ChatID)
		if err != nil {
			return Response{}, fmt.Errorf("getting clientID by chatID %s: %v", p.ChatID.String(), err)
		}
		res.Chats = append(res.Chats, &Chat{ChatID: p.ChatID, ClientID: clientID})
	}

	return res, nil
}
