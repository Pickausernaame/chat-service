package managerevents_test

import (
	"encoding/json"
	"testing"

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
