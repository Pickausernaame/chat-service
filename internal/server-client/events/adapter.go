package clientevents

import (
	"errors"

	eventstream "github.com/Pickausernaame/chat-service/internal/services/event-stream"
	websocketstream "github.com/Pickausernaame/chat-service/internal/websocket-stream"
	"github.com/Pickausernaame/chat-service/pkg/pointer"
)

var _ websocketstream.EventAdapter = Adapter{}

var ErrInvalidEventType = errors.New("invalid type")

type Adapter struct{}

func (Adapter) Adapt(ev eventstream.Event) (any, error) {
	switch ev.Type() {
	case eventstream.EventTypeNewMessageEvent:
		res, ok := ev.(*eventstream.NewMessageEvent)
		if !ok {
			return nil, ErrInvalidEventType
		}
		return toNewMessageEvent(res), nil
	case eventstream.EventTypeMessageSentEvent:
		res, ok := ev.(*eventstream.MessageSentEvent)
		if !ok {
			return nil, ErrInvalidEventType
		}
		return toMessageSentEvent(res), nil
	case eventstream.EventTypeMessageBlockedEvent:
		res, ok := ev.(*eventstream.MessageBlockedEvent)
		if !ok {
			return nil, ErrInvalidEventType
		}
		return toMessageBlockedEvent(res), nil
	}
	return nil, errors.New("invalid event")
}

func toNewMessageEvent(ev *eventstream.NewMessageEvent) *NewMessageEvent {
	return &NewMessageEvent{
		AuthorId:  pointer.PtrWithZeroAsNil(ev.UserID),
		Body:      ev.MessageBody,
		CreatedAt: ev.CreatedAt,
		EventId:   ev.EventID,
		EventType: NewMessageEventEventTypeNewMessageEvent,
		IsService: ev.IsService,
		MessageId: ev.MessageID,
		RequestId: ev.RequestID,
	}
}

func toMessageSentEvent(ev *eventstream.MessageSentEvent) *MessageSentEvent {
	return &MessageSentEvent{
		EventId:   ev.EventID,
		EventType: MessageSentEventEventTypeMessageSentEvent,
		MessageId: ev.MessageID,
		RequestId: ev.RequestID,
	}
}

func toMessageBlockedEvent(ev *eventstream.MessageBlockedEvent) *MessageBlockedEvent {
	return &MessageBlockedEvent{
		EventId:   ev.EventID,
		EventType: MessageBlockedEventEventTypeMessageBlockedEvent,
		MessageId: ev.MessageID,
		RequestId: ev.RequestID,
	}
}
