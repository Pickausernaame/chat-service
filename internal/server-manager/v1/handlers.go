package managerv1

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	canreceiveproblems "github.com/Pickausernaame/chat-service/internal/usecases/manager/can-receive-problems"
	getassignedproblems "github.com/Pickausernaame/chat-service/internal/usecases/manager/get-assigned-problems"
	getchathistory "github.com/Pickausernaame/chat-service/internal/usecases/manager/get-chat-history"
	sendmessage "github.com/Pickausernaame/chat-service/internal/usecases/manager/send-message"
	setreadyreceiveproblems "github.com/Pickausernaame/chat-service/internal/usecases/manager/set-ready-receive-problems"
)

var _ ServerInterface = (*Handlers)(nil)

//go:generate mockgen -source=$GOFILE -destination=mocks/handlers_mocks.gen.go -package=managerv1mocks

type sendMessageUseCase interface {
	Handle(ctx context.Context, req sendmessage.Request) (sendmessage.Response, error)
}

type canReceiveProblemsUseCase interface {
	Handle(ctx context.Context, req canreceiveproblems.Request) (canreceiveproblems.Response, error)
}

type setReadyReceiveProblemsUseCase interface {
	Handle(ctx context.Context, req setreadyreceiveproblems.Request) (setreadyreceiveproblems.Response, error)
}

type getAssignedProblemsUseCase interface {
	Handle(ctx context.Context, req getassignedproblems.Request) (getassignedproblems.Response, error)
}

type getChatHistoryUseCase interface {
	Handle(ctx context.Context, req getchathistory.Request) (getchathistory.Response, error)
}

//go:generate options-gen -out-filename=handlers.gen.go -from-struct=Options
type Options struct {
	canReceiveProblemsUseCase      canReceiveProblemsUseCase      `option:"mandatory" validate:"required"`
	setReadyReceiveProblemsUseCase setReadyReceiveProblemsUseCase `option:"mandatory" validate:"required"`
	getAssignedProblemsUseCase     getAssignedProblemsUseCase     `option:"mandatory" validate:"required"`
	getChatHistoryUseCase          getChatHistoryUseCase          `option:"mandatory" validate:"required"`
	sendMessageUseCase             sendMessageUseCase             `option:"mandatory" validate:"required"`
}

type Handlers struct {
	Options
	lg *zap.Logger
}

func NewHandlers(opts Options) (Handlers, error) {
	if err := opts.Validate(); err != nil {
		return Handlers{}, fmt.Errorf("validate options: %v", err)
	}
	return Handlers{Options: opts, lg: zap.L().Named("managerv1-handler")}, nil
}
