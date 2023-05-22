package dummy

import (
	"context"

	eventstream "github.com/Pickausernaame/chat-service/internal/services/event-stream"
	"github.com/Pickausernaame/chat-service/internal/types"
)

type DummyEventStream struct{}

func (DummyEventStream) Subscribe(ctx context.Context, _ types.UserID) (<-chan eventstream.Event, error) {
	events := make(chan eventstream.Event)
	go func() {
		defer close(events)
		<-ctx.Done()
	}()
	return events, nil
}
