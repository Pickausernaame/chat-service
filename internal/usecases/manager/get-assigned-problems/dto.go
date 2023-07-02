package getassignedproblems

import (
	"github.com/Pickausernaame/chat-service/internal/types"
	"github.com/Pickausernaame/chat-service/internal/validator"
)

type Request struct {
	ManagerID types.UserID `validate:"required"`
}

func (r Request) Validate() error {
	return validator.Validator.Struct(r)
}

type Response struct {
	Chats []*Chat
}

type Chat struct {
	ChatID   types.ChatID
	ClientID types.UserID
}
