package messagesrepo

import (
	"time"

	"github.com/Pickausernaame/chat-service/internal/store"
	"github.com/Pickausernaame/chat-service/internal/types"
)

type Message struct {
	ID                  types.MessageID
	ChatID              types.ChatID
	AuthorID            types.UserID
	InitialRequestID    types.RequestID
	Body                string
	IsVisibleForClient  bool
	IsVisibleForManager bool
	IsBlocked           bool
	IsService           bool
	CreatedAt           time.Time
}

func adaptStoreMessage(m *store.Message) *Message {
	return &Message{
		ID:                  m.ID,
		ChatID:              m.ChatID,
		AuthorID:            m.AuthorID,
		InitialRequestID:    m.InitialRequestID,
		Body:                m.Body,
		IsBlocked:           m.IsBlocked,
		IsVisibleForClient:  m.IsVisibleForClient,
		IsVisibleForManager: m.IsVisibleForManager,
		IsService:           m.IsService,
		CreatedAt:           m.CreatedAt,
	}
}
