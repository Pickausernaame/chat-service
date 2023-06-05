package managerassignedtoproblemjob_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	managerassignedtoproblemjob "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/manager-assigned-to-problem"
	"github.com/Pickausernaame/chat-service/internal/types"
)

func TestMarshalPayload_Smoke(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		p, err := managerassignedtoproblemjob.MarshalPayload(
			types.NewUserID(),
			types.NewUserID(),
			types.NewRequestID(),
			types.NewMessageID())
		require.NoError(t, err)
		assert.NotEmpty(t, p)
	})

	t.Run("invalid input", func(t *testing.T) {
		p, err := managerassignedtoproblemjob.MarshalPayload(
			types.UserIDNil,
			types.NewUserID(),
			types.NewRequestID(),
			types.NewMessageID())
		require.Error(t, err)
		assert.Empty(t, p)

		p, err = managerassignedtoproblemjob.MarshalPayload(
			types.NewUserID(),
			types.UserIDNil,
			types.NewRequestID(),
			types.NewMessageID())
		require.Error(t, err)
		assert.Empty(t, p)

		p, err = managerassignedtoproblemjob.MarshalPayload(
			types.NewUserID(),
			types.NewUserID(),
			types.RequestIDNil,
			types.NewMessageID())
		require.Error(t, err)
		assert.Empty(t, p)

		p, err = managerassignedtoproblemjob.MarshalPayload(
			types.NewUserID(),
			types.NewUserID(),
			types.NewRequestID(),
			types.MessageIDNil)
		require.Error(t, err)
		assert.Empty(t, p)
	})
}
