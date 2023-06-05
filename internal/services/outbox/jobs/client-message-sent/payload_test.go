package clientmessagesentjob_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	clientmessagesentjob "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/client-message-sent"
	"github.com/Pickausernaame/chat-service/internal/types"
)

func TestMarshalPayload_Smoke(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		p, err := clientmessagesentjob.MarshalPayload(types.NewMessageID())
		require.NoError(t, err)
		assert.NotEmpty(t, p)
	})

	t.Run("invalid input", func(t *testing.T) {
		p, err := clientmessagesentjob.MarshalPayload(types.MessageIDNil)
		require.Error(t, err)
		assert.Empty(t, p)
	})
}
