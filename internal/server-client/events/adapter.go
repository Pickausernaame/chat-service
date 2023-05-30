package clientevents

import (
	"errors"

	eventstream "github.com/Pickausernaame/chat-service/internal/services/event-stream"
	websocketstream "github.com/Pickausernaame/chat-service/internal/websocket-stream"
	"github.com/Pickausernaame/chat-service/pkg/pointer"
)

var _ websocketstream.EventAdapter = Adapter{}

type Adapter struct{}

func (Adapter) Adapt(ev eventstream.Event) (any, error) {
	switch ev.Type() {
	case eventstream.EventTypeNewMessageEvent:
		res := ev.(*eventstream.NewMessageEvent)
		return toNewMessageEvent(res), nil
	case eventstream.EventTypeMessageSentEvent:
		res := ev.(*eventstream.MessageSentEvent)
		return toMessageSentEvent(res), nil
	case eventstream.EventTypeMessageBlockedEvent:
		res := ev.(*eventstream.MessageBlockedEvent)
		return toMessageBlockedEvent(res), nil
	}
	return nil, errors.New("invalid event")
}

func toNewMessageEvent(ev *eventstream.NewMessageEvent) *NewMessageEvent {
	return &NewMessageEvent{
		AuthorId:  pointer.PtrWithZeroAsNil(ev.UserID),
		Body:      pointer.PtrWithZeroAsNil(ev.MessageBody),
		CreatedAt: pointer.PtrWithZeroAsNil(ev.CreatedAt),
		EventId:   pointer.PtrWithZeroAsNil(ev.EventID),
		EventType: pointer.PtrWithZeroAsNil(NewMessageEventEventTypeNewMessageEvent),
		IsService: pointer.PtrWithZeroAsNil(ev.IsService),
		MessageId: pointer.PtrWithZeroAsNil(ev.MessageID),
		RequestId: pointer.PtrWithZeroAsNil(ev.RequestID),
	}
}

func toMessageSentEvent(ev *eventstream.MessageSentEvent) *MessageSentEvent {
	return &MessageSentEvent{
		EventId:   pointer.PtrWithZeroAsNil(ev.EventID),
		EventType: pointer.PtrWithZeroAsNil(BaseEventEventTypeMessageSentEvent),
		MessageId: pointer.PtrWithZeroAsNil(ev.MessageID),
		RequestId: pointer.PtrWithZeroAsNil(ev.RequestID),
	}
}

func toMessageBlockedEvent(ev *eventstream.MessageBlockedEvent) *MessageBlockedEvent {
	return &MessageBlockedEvent{
		EventId:   pointer.PtrWithZeroAsNil(ev.EventID),
		EventType: pointer.PtrWithZeroAsNil(BaseEventEventTypeMessageBlockedEvent),
		MessageId: pointer.PtrWithZeroAsNil(ev.MessageID),
		RequestId: pointer.PtrWithZeroAsNil(ev.RequestID),
	}
}
