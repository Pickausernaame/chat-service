package websocketstream

import (
	"encoding/json"
	"io"

	eventstream "github.com/Pickausernaame/chat-service/internal/services/event-stream"
)

// EventAdapter converts the event from the stream to the appropriate object.
type EventAdapter interface {
	Adapt(event eventstream.Event) (any, error)
}

// EventWriter write adapted event it to the socket.
type EventWriter interface {
	Write(event any, out io.Writer) error
}

type JSONEventWriter struct{}

func (JSONEventWriter) Write(event any, out io.Writer) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, err = out.Write(payload)
	if err != nil {
		return err
	}
	return nil
}

type DummyAdapter struct{}

func (DummyAdapter) Adapt(event eventstream.Event) (any, error) {
	return event, nil
}
