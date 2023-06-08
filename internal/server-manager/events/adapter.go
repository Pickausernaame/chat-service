package managerevents

import (
	"errors"

	eventstream "github.com/Pickausernaame/chat-service/internal/services/event-stream"
	websocketstream "github.com/Pickausernaame/chat-service/internal/websocket-stream"
)

var ErrInvalidEventType = errors.New("invalid type")

var _ websocketstream.EventAdapter = Adapter{}

type Adapter struct{}

func (Adapter) Adapt(ev eventstream.Event) (any, error) {
	switch ev.Type() {
	case eventstream.EventTypeNewChatEvent:
		res, ok := ev.(*eventstream.NewChatEvent)
		if !ok {
			return nil, ErrInvalidEventType
		}
		return toNewChatEvent(res), nil
	case eventstream.EventTypeNewMessageEvent:
		res, ok := ev.(*eventstream.NewMessageEvent)
		if !ok {
			return nil, ErrInvalidEventType
		}
		return toNewMessageEvent(res), nil
	case eventstream.EventTypeChatClosedEvent:
		res, ok := ev.(*eventstream.ChatClosedEvent)
		if !ok {
			return nil, ErrInvalidEventType
		}
		return toChatClosedEvent(res), nil
	}
	return nil, nil
}

func toNewChatEvent(ev *eventstream.NewChatEvent) *NewChatEvent {
	return &NewChatEvent{
		CanTakeMoreProblems: ev.CanTakeMoreProblems,
		ChatId:              ev.ChatID,
		ClientId:            ev.ClientID,
		EventId:             ev.EventID,
		EventType:           NewChatEventEventTypeNewChatEvent,
		RequestId:           ev.RequestID,
	}
}

func toNewMessageEvent(ev *eventstream.NewMessageEvent) *NewMessageEvent {
	return &NewMessageEvent{
		AuthorId:  ev.UserID,
		Body:      ev.MessageBody,
		CreatedAt: ev.CreatedAt,
		EventId:   ev.EventID,
		ChatId:    ev.ChatID,
		EventType: NewMessageEventEventTypeNewMessageEvent,
		MessageId: ev.MessageID,
		RequestId: ev.RequestID,
	}
}

func toChatClosedEvent(ev *eventstream.ChatClosedEvent) *ChatClosedEvent {
	return &ChatClosedEvent{
		EventId:             ev.EventID,
		ChatId:              ev.ChatID,
		EventType:           ChatClosedEventEventTypeNewChatEvent,
		RequestId:           ev.RequestID,
		CanTakeMoreProblems: ev.CanTakeMoreProblems,
	}
}
