package managerload

import (
	"context"

	"go.uber.org/zap"

	"github.com/Pickausernaame/chat-service/internal/types"
)

const (
	serviceName = "manager-pool"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/service_mock.gen.go -package=managerpool_mock

type problemsRepository interface {
	GetManagerOpenProblemsCount(ctx context.Context, managerID types.UserID) (int, error)
}

//go:generate options-gen -out-filename=service_options.gen.go -from-struct=Options
type Options struct {
	maxProblemsAtTime int `option:"mandatory" validate:"required,min=1,max=30"`

	problemsRepo problemsRepository `option:"mandatory" validate:"required"`
}

type Service struct {
	Options
	lg *zap.Logger
}

func New(opts Options) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return &Service{Options: opts, lg: zap.L().Named(serviceName)}, nil
}
