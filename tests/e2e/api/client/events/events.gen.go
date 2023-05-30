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

// Defines values for NewMessageEventEventType.
const (
	NewMessageEventEventTypeMessageBlockedEvent NewMessageEventEventType = "MessageBlockedEvent"
	NewMessageEventEventTypeMessageSentEvent    NewMessageEventEventType = "MessageSentEvent"
	NewMessageEventEventTypeNewMessageEvent     NewMessageEventEventType = "NewMessageEvent"
)

// BaseEvent defines model for BaseEvent.
type BaseEvent struct {
	// EventId Unique identifier for the event
	EventId *types.EventID `json:"eventId,omitempty"`

	// EventType Type of the event
	EventType *BaseEventEventType `json:"eventType,omitempty"`

	// MessageId Unique identifier for the message
	MessageId *types.MessageID `json:"messageId,omitempty"`

	// RequestId Unique identifier for the request
	RequestId *types.RequestID `json:"requestId,omitempty"`
}

// BaseEventEventType Type of the event
type BaseEventEventType string

// MessageBlockedEvent defines model for MessageBlockedEvent.
type MessageBlockedEvent = BaseEvent

// MessageSentEvent defines model for MessageSentEvent.
type MessageSentEvent = BaseEvent

// NewMessageEvent defines model for NewMessageEvent.
type NewMessageEvent struct {
	// AuthorId Unique identifier for the author
	AuthorId *types.UserID `json:"authorId,omitempty"`

	// Body Body of the message
	Body *string `json:"body,omitempty"`

	// CreatedAt Date and time of event creation
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	// EventId Unique identifier for the event
	EventId *types.EventID `json:"eventId,omitempty"`

	// EventType Type of the event
	EventType *NewMessageEventEventType `json:"eventType,omitempty"`

	// IsService Indicates if the event is a service event
	IsService *bool `json:"isService,omitempty"`

	// MessageId Unique identifier for the message
	MessageId *types.MessageID `json:"messageId,omitempty"`

	// RequestId Unique identifier for the request
	RequestId *types.RequestID `json:"requestId,omitempty"`
}

// NewMessageEventEventType Type of the event
type NewMessageEventEventType string
