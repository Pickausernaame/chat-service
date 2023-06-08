package afcverdictsprocessor

import (
	"github.com/Pickausernaame/chat-service/internal/types"
)

type verdict struct {
	ChatID    types.ChatID    `json:"chatId" validate:"required"`
	MessageID types.MessageID `json:"messageId" validate:"required"`
	Status    string          `json:"status" validate:"required,oneof=suspicious ok"`
}
