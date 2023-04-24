package clientv1

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/Pickausernaame/chat-service/internal/validator"
)

//go:generate options-gen -out-filename=handler_options.gen.go -from-struct=Options
type Options struct {
	logger *zap.Logger `option:"mandatory" validate:"required"`
	// Ждут своего часа.
}

type Handlers struct {
	Options
}

func NewHandlers(opts Options) (Handlers, error) {
	if err := validator.Validator.Struct(opts); err != nil {
		return Handlers{}, fmt.Errorf("options validation error: %v", err)
	}
	return Handlers{Options: opts}, nil
}
