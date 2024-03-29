// Package apiclientevents provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.5-0.20230506011706-29ebe3262399 DO NOT EDIT.
package apiclientevents

import (
	"time"

	"github.com/Pickausernaame/chat-service/internal/types"
)

// Defines values for BaseEventEventType.
const (
	BaseEventEventTypeMessageBlockedEvent BaseEventEventType = "MessageBlockedEvent"
	BaseEventEventTypeMessageSentEvent    BaseEventEventType = "MessageSentEvent"
	BaseEventEventTypeNewMessageEvent     BaseEventEventType = "NewMessageEvent"
)

// Defines values for MessageBlockedEventEventType.
const (
	MessageBlockedEventEventTypeMessageBlockedEvent MessageBlockedEventEventType = "MessageBlockedEvent"
	MessageBlockedEventEventTypeMessageSentEvent    MessageBlockedEventEventType = "MessageSentEvent"
	MessageBlockedEventEventTypeNewMessageEvent     MessageBlockedEventEventType = "NewMessageEvent"
)

// Defines values for MessageSentEventEventType.
const (
	MessageSentEventEventTypeMessageBlockedEvent MessageSentEventEventType = "MessageBlockedEvent"
	MessageSentEventEventTypeMessageSentEvent    MessageSentEventEventType = "MessageSentEvent"
	MessageSentEventEventTypeNewMessageEvent     MessageSentEventEventType = "NewMessageEvent"
)

// Defines values for NewMessageEventEventType.
const (
	NewMessageEventEventTypeMessageBlockedEvent NewMessageEventEventType = "MessageBlockedEvent"
	NewMessageEventEventTypeMessageSentEvent    NewMessageEventEventType = "MessageSentEvent"
	NewMessageEventEventTypeNewMessageEvent     NewMessageEventEventType = "NewMessageEvent"
)

// BaseEvent defines model for BaseEvent.
type BaseEvent struct {
	// EventId Unique identifier for the event
	EventId types.EventID `json:"eventId"`

	// EventType Type of the event
	EventType BaseEventEventType `json:"eventType"`

	// RequestId Unique identifier for the request
	RequestId types.RequestID `json:"requestId"`
}

// BaseEventEventType Type of the event
type BaseEventEventType string

// MessageBlockedEvent defines model for MessageBlockedEvent.
type MessageBlockedEvent struct {
	// EventId Unique identifier for the event
	EventId types.EventID `json:"eventId"`

	// EventType Type of the event
	EventType MessageBlockedEventEventType `json:"eventType"`
	MessageId types.MessageID              `json:"messageId"`

	// RequestId Unique identifier for the request
	RequestId types.RequestID `json:"requestId"`
}

// MessageBlockedEventEventType Type of the event
type MessageBlockedEventEventType string

// MessageId defines model for MessageId.
type MessageId struct {
	MessageId types.MessageID `json:"messageId"`
}

// MessageSentEvent defines model for MessageSentEvent.
type MessageSentEvent struct {
	// EventId Unique identifier for the event
	EventId types.EventID `json:"eventId"`

	// EventType Type of the event
	EventType MessageSentEventEventType `json:"eventType"`
	MessageId types.MessageID           `json:"messageId"`

	// RequestId Unique identifier for the request
	RequestId types.RequestID `json:"requestId"`
}

// MessageSentEventEventType Type of the event
type MessageSentEventEventType string

// NewMessageEvent defines model for NewMessageEvent.
type NewMessageEvent struct {
	// AuthorId Unique identifier for the author
	AuthorId *types.UserID `json:"authorId,omitempty"`

	// Body Body of the message
	Body string `json:"body"`

	// CreatedAt Date and time of event creation
	CreatedAt time.Time `json:"createdAt"`

	// EventId Unique identifier for the event
	EventId types.EventID `json:"eventId"`

	// EventType Type of the event
	EventType NewMessageEventEventType `json:"eventType"`

	// IsService Indicates if the event is a service event
	IsService bool            `json:"isService"`
	MessageId types.MessageID `json:"messageId"`

	// RequestId Unique identifier for the request
	RequestId types.RequestID `json:"requestId"`
}

// NewMessageEventEventType Type of the event
type NewMessageEventEventType string
