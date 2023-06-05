package eventstream

import (
	"fmt"
)

// Matches checks if the actual event matches the expected event, ignoring the EventID field.
func (m *NewMessageEvent) Matches(x interface{}) bool {
	actualEvent, ok := x.(*NewMessageEvent)
	if !ok {
		return false
	}

	// Compare the fields excluding EventID
	return actualEvent.IsService == m.IsService &&
		actualEvent.MessageBody == m.MessageBody &&
		actualEvent.EventType == m.EventType &&
		actualEvent.MessageID == m.MessageID &&
		actualEvent.CreatedAt == m.CreatedAt &&
		actualEvent.ChatID == m.ChatID &&
		actualEvent.RequestID == m.RequestID
}

func (m *NewMessageEvent) String() string {
	return fmt.Sprintf("%v", *m)
}

func (m *MessageBlockedEvent) Matches(x interface{}) bool {
	actualEvent, ok := x.(*MessageBlockedEvent)
	if !ok {
		return false
	}

	// Compare the fields excluding EventID
	return actualEvent.EventType == m.EventType &&
		actualEvent.MessageID == m.MessageID &&
		actualEvent.RequestID == m.RequestID
}

func (m *MessageBlockedEvent) String() string {
	return fmt.Sprintf("%v", *m)
}

func (m *MessageSentEvent) Matches(x interface{}) bool {
	actualEvent, ok := x.(*MessageSentEvent)
	if !ok {
		return false
	}

	// Compare the fields excluding EventID
	return actualEvent.EventType == m.EventType &&
		actualEvent.MessageID == m.MessageID &&
		actualEvent.RequestID == m.RequestID
}

func (m *MessageSentEvent) String() string {
	return fmt.Sprintf("%v", *m)
}

func (m *NewChatEvent) Matches(x interface{}) bool {
	actualEvent, ok := x.(*NewChatEvent)
	if !ok {
		return false
	}

	// Compare the fields excluding EventID
	return actualEvent.EventType == m.EventType &&
		actualEvent.RequestID == m.RequestID &&
		actualEvent.ChatID == m.ChatID &&
		actualEvent.ClientID == m.ClientID
}

func (m *NewChatEvent) String() string {
	return fmt.Sprintf("%v", *m)
}
