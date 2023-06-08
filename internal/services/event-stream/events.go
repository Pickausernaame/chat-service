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
	EventTypeNewChatEvent        string = "NewChatEvent"
	EventTypeChatClosedEvent     string = "ChatClosedEvent"
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
	UserID      types.UserID
	CreatedAt   time.Time `validate:"required"`
	MessageBody string    `validate:"required"`
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

func (m *NewMessageEvent) Validate() error {
	return validator.Validator.Struct(m)
}

func (m *NewMessageEvent) Type() string {
	return m.EventType
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

type MessageBlockedEvent struct {
	event
	EventType string
	EventID   types.EventID   `validate:"required"`
	RequestID types.RequestID `validate:"required"`
	MessageID types.MessageID `validate:"required"`
}

func NewMessageBlockedEvent(
	eventID types.EventID,
	requestID types.RequestID,
	messageID types.MessageID,
) *MessageBlockedEvent {
	return &MessageBlockedEvent{
		event:     event{},
		EventID:   eventID,
		RequestID: requestID,
		MessageID: messageID,
		EventType: EventTypeMessageBlockedEvent,
	}
}

func (e *MessageBlockedEvent) Validate() error {
	return validator.Validator.Struct(e)
}

func (e *MessageBlockedEvent) Type() string {
	return e.EventType
}

type NewChatEvent struct {
	event

	EventID             types.EventID   `validate:"required"`
	RequestID           types.RequestID `validate:"required"`
	CanTakeMoreProblems bool
	ChatID              types.ChatID `validate:"required"`
	ClientID            types.UserID `validate:"required"`
	EventType           string       `validate:"required"`
}

func NewNewChatEvent(
	eventID types.EventID,
	requestID types.RequestID,
	canTakeMoreProblems bool,
	chatID types.ChatID,
	clientID types.UserID,
) *NewChatEvent {
	return &NewChatEvent{
		event:               event{},
		EventID:             eventID,
		RequestID:           requestID,
		CanTakeMoreProblems: canTakeMoreProblems,
		ChatID:              chatID,
		ClientID:            clientID,
		EventType:           EventTypeNewChatEvent,
	}
}

func (e *NewChatEvent) Validate() error {
	return validator.Validator.Struct(e)
}

func (e *NewChatEvent) Type() string {
	return e.EventType
}

type ChatClosedEvent struct {
	event

	EventID             types.EventID   `validate:"required"`
	RequestID           types.RequestID `validate:"required"`
	CanTakeMoreProblems bool
	ChatID              types.ChatID `validate:"required"`
	EventType           string       `validate:"required"`
}

func NewChatClosedEvent(
	eventID types.EventID,
	requestID types.RequestID,
	canTakeMoreProblems bool,
	chatID types.ChatID,
) *ChatClosedEvent {
	return &ChatClosedEvent{
		event:               event{},
		EventID:             eventID,
		RequestID:           requestID,
		CanTakeMoreProblems: canTakeMoreProblems,
		ChatID:              chatID,
		EventType:           EventTypeChatClosedEvent,
	}
}

func (e *ChatClosedEvent) Validate() error {
	return validator.Validator.Struct(e)
}

func (e *ChatClosedEvent) Type() string {
	return e.EventType
}
