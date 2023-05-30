package setreadyreceiveproblems

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/Pickausernaame/chat-service/internal/types"
)

var (
	ErrInvalidRequest  = errors.New("invalid request")
	ErrManagerOverload = errors.New("manager overload")
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=setreadyreceiveproblemsmocks

type managerLoadService interface {
	CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error)
}

type managerPool interface {
	Put(ctx context.Context, managerID types.UserID) error
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	managerLoadService managerLoadService `option:"mandatory" validate:"required"`
	managerPool        managerPool        `option:"mandatory" validate:"required"`
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
		lg:      zap.L().Named("set-ready-receive-problems-usecase"),
	}, nil
}

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if err := req.Validate(); err != nil {
		return Response{}, ErrInvalidRequest
	}

	ok, err := u.managerLoadService.CanManagerTakeProblem(ctx, req.ManagerID)
	if err != nil {
		return Response{}, fmt.Errorf("checking can manager take problem: %v", err)
	}
	if !ok {
		return Response{}, ErrManagerOverload
	}

	if err = u.managerPool.Put(ctx, req.ManagerID); err != nil {
		return Response{}, fmt.Errorf("putting manager to pool: %v", err)
	}

	return Response{}, nil
}
