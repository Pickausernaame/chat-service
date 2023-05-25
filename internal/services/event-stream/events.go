package eventstream

import (
	"time"

	"github.com/Pickausernaame/chat-service/internal/types"
	"github.com/Pickausernaame/chat-service/internal/validator"
)

type Event interface {
	eventMarker()
	Validate() error
	Type() string
}

const (
	EventTypeMessageBlockedEvent string = "MessageBlockedEvent"
	EventTypeMessageSentEvent    string = "MessageSentEvent"
	EventTypeNewMessageEvent     string = "NewMessageEvent"
)

type event struct{}         //
func (*event) eventMarker() {}

// NewMessageEvent is a signal about the appearance of a new message in the chat.
type NewMessageEvent struct {
	event
	EventType   string
	EventID     types.EventID   `validate:"required"`
	RequestID   types.RequestID `validate:"required"`
	ChatID      types.ChatID    `validate:"required"`
	MessageID   types.MessageID `validate:"required"`
	UserID      types.UserID    `validate:"required"`
	CreatedAt   time.Time       `validate:"required"`
	MessageBody string          `validate:"required"`
	IsService   bool
}

func NewNewMessageEvent(eventID types.EventID,
	requestID types.RequestID, chatID types.ChatID,
	messageID types.MessageID, userID types.UserID,
	createdAt time.Time, body string, isService bool,
) *NewMessageEvent {
	return &NewMessageEvent{
		event:       event{},
		EventID:     eventID,
		RequestID:   requestID,
		ChatID:      chatID,
		MessageID:   messageID,
		UserID:      userID,
		CreatedAt:   createdAt,
		MessageBody: body,
		IsService:   isService,
		EventType:   EventTypeNewMessageEvent,
	}
}

func (e *NewMessageEvent) Validate() error {
	return validator.Validator.Struct(e)
}

func (e *NewMessageEvent) Type() string {
	return e.EventType
}

type MessageSentEvent struct {
	event
	EventType string
	EventID   types.EventID   `validate:"required"`
	RequestID types.RequestID `validate:"required"`
	MessageID types.MessageID `validate:"required"`
}

func NewMessageSentEvent(
	eventID types.EventID,
	requestID types.RequestID,
	messageID types.MessageID,
) *MessageSentEvent {
	return &MessageSentEvent{
		event:     event{},
		EventID:   eventID,
		RequestID: requestID,
		MessageID: messageID,
		EventType: EventTypeMessageSentEvent,
	}
}

func (e *MessageSentEvent) Validate() error {
	return validator.Validator.Struct(e)
}

func (e *MessageSentEvent) Type() string {
	return e.EventType
}
