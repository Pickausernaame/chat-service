package sendmanagermessagejob_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sendmanagermessagejob "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/send-manager-message"
	"github.com/Pickausernaame/chat-service/internal/types"
)

func TestMarshalPayload_Smoke(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		p, err := sendmanagermessagejob.MarshalPayload(types.NewMessageID())
		require.NoError(t, err)
		assert.NotEmpty(t, p)
	})

	t.Run("invalid input", func(t *testing.T) {
		p, err := sendmanagermessagejob.MarshalPayload(types.MessageIDNil)
		require.Error(t, err)
		assert.Empty(t, p)
	})
}
