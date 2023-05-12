package gethistory

import (
	"fmt"
	"time"

	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	"github.com/Pickausernaame/chat-service/internal/types"
	"github.com/Pickausernaame/chat-service/internal/validator"
)

type Request struct {
	ID       types.RequestID `validate:"required"`
	ClientID types.UserID    `validate:"required"`
	PageSize int             `validate:"omitempty,gte=10,lte=100"`
	Cursor   string          `validate:"omitempty,base64url"`
}

func (r Request) Validate() error {
	if r.PageSize == 0 && r.Cursor == "" {
		return fmt.Errorf("page size or cursor not specified")
	}

	if r.PageSize != 0 && r.Cursor != "" {
		return fmt.Errorf("page size AND cursor specified")
	}
	return validator.Validator.Struct(r)
}

type Response struct {
	Messages   []*Message
	NextCursor string
}

type Message struct {
	ID         types.MessageID
	AuthorID   types.UserID
	Body       string
	CreatedAt  time.Time
	IsBlocked  bool
	IsReceived bool
	IsService  bool
}

func toDTOMessage(msg *messagesrepo.Message) *Message {
	return &Message{
		ID:         msg.ID,
		AuthorID:   msg.AuthorID,
		Body:       msg.Body,
		CreatedAt:  msg.CreatedAt,
		IsBlocked:  msg.IsBlocked,
		IsService:  msg.IsService,
		IsReceived: msg.IsVisibleForManager && !msg.IsBlocked,
	}
}
