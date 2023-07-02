package managerevents_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	managerevents "github.com/Pickausernaame/chat-service/internal/server-manager/events"
	eventstream "github.com/Pickausernaame/chat-service/internal/services/event-stream"
	"github.com/Pickausernaame/chat-service/internal/types"
)

func TestAdapter_Adapt(t *testing.T) {
	cases := []struct {
		name    string
		ev      eventstream.Event
		expJSON string
	}{
		{
			name: "smoke",
			ev: eventstream.NewNewChatEvent(
				types.MustParse[types.EventID]("d0ffbd36-bc30-11ed-8286-461e464ebed8"),
				types.MustParse[types.RequestID]("cee5f290-bc30-11ed-b7fe-461e464ebed8"),
				true,
				types.MustParse[types.ChatID]("cb36a888-bc30-11ed-b843-461e464ebed8"),
				types.MustParse[types.UserID]("eb36a888-bc30-11ed-b843-461e464ebed8"),
			),
			expJSON: `{
				"eventId": "d0ffbd36-bc30-11ed-8286-461e464ebed8",
				"requestId": "cee5f290-bc30-11ed-b7fe-461e464ebed8",
				"canTakeMoreProblems": true,
				"chatId": "cb36a888-bc30-11ed-b843-461e464ebed8",
				"clientId": "eb36a888-bc30-11ed-b843-461e464ebed8",
				"eventType": "NewChatEvent"
			}`,
		},

		{
			name: "service message",
			ev: eventstream.NewNewMessageEvent(
				types.MustParse[types.EventID]("d0ffbd36-bc30-11ed-8286-461e464ebed8"),
				types.MustParse[types.RequestID]("cee5f290-bc30-11ed-b7fe-461e464ebed8"),
				types.MustParse[types.ChatID]("31b4dc06-bc31-11ed-93cc-461e464ebed8"),
				types.MustParse[types.MessageID]("cb36a888-bc30-11ed-b843-461e464ebed8"),
				types.MustParse[types.UserID]("7d67b14d-221e-4499-9be2-6707d7df1adc"),
				time.Unix(1, 1).UTC(),
				"Manager will coming soon",
				true,
			),
			expJSON: `{
				"authorId": "7d67b14d-221e-4499-9be2-6707d7df1adc",
				"body": "Manager will coming soon",
				"chatId": "31b4dc06-bc31-11ed-93cc-461e464ebed8",
				"createdAt": "1970-01-01T00:00:01.000000001Z",
				"eventId": "d0ffbd36-bc30-11ed-8286-461e464ebed8",
				"eventType": "NewMessageEvent",
				"messageId": "cb36a888-bc30-11ed-b843-461e464ebed8",
				"requestId": "cee5f290-bc30-11ed-b7fe-461e464ebed8"
			}`,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			adapted, err := managerevents.Adapter{}.Adapt(tt.ev)
			require.NoError(t, err)

			raw, err := json.Marshal(adapted)
			require.NoError(t, err)
			assert.JSONEq(t, tt.expJSON, string(raw))
		})
	}
}
