package eventstream

import (
	"time"

	"github.com/Pickausernaame/chat-service/internal/types"
	"github.com/Pickausernaame/chat-service/internal/validator"
)

type Event interface {
	eventMarker()
	Validate() error
}

type event struct{}         //
func (*event) eventMarker() {}

// NewMessageEvent is a signal about the appearance of a new message in the chat.
type NewMessageEvent struct {
	event
	eventID     types.EventID   `validate:"required"`
	requestID   types.RequestID `validate:"required"`
	chatID      types.ChatID    `validate:"required"`
	messageID   types.MessageID `validate:"required"`
	userID      types.UserID    `validate:"required"`
	createdAt   time.Time       `validate:"required"`
	MessageBody string          `validate:"required"`
}

func NewNewMessageEvent(eventID types.EventID,
	requestID types.RequestID, chatID types.ChatID,
	messageID types.MessageID, userID types.UserID, createdAt time.Time, body string, _ bool) *NewMessageEvent {
	return &NewMessageEvent{
		event:       event{},
		eventID:     eventID,
		requestID:   requestID,
		chatID:      chatID,
		messageID:   messageID,
		userID:      userID,
		createdAt:   createdAt,
		MessageBody: body,
	}
}

func (e *NewMessageEvent) Validate() error {
	return validator.Validator.Struct(e)
}
