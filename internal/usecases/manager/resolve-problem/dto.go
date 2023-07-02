package resolveproblem

import (
	"github.com/Pickausernaame/chat-service/internal/types"
	"github.com/Pickausernaame/chat-service/internal/validator"
)

type Request struct {
	ChatID    types.ChatID    `validate:"required"`
	ManagerID types.UserID    `validate:"required"`
	RequestID types.RequestID `validate:"required"`
}

func (r Request) Validate() error {
	return validator.Validator.Struct(r)
}

type Response struct{}
