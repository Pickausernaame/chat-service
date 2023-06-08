package jobresolveproblem_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	jobresolveproblem "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/job-resolve-problem"
	"github.com/Pickausernaame/chat-service/internal/types"
)

func TestMarshalPayload_Smoke(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		p, err := jobresolveproblem.MarshalPayload(
			types.NewUserID(),
			types.NewRequestID(),
			types.NewMessageID(),
			types.NewChatID())
		require.NoError(t, err)
		assert.NotEmpty(t, p)
	})

	t.Run("invalid input", func(t *testing.T) {
		p, err := jobresolveproblem.MarshalPayload(
			types.UserIDNil,
			types.NewRequestID(),
			types.NewMessageID(),
			types.NewChatID())
		require.Error(t, err)
		assert.Empty(t, p)

		p, err = jobresolveproblem.MarshalPayload(
			types.NewUserID(),
			types.RequestIDNil,
			types.NewMessageID(),
			types.NewChatID())
		require.Error(t, err)
		assert.Empty(t, p)

		p, err = jobresolveproblem.MarshalPayload(
			types.NewUserID(),
			types.NewRequestID(),
			types.MessageIDNil,
			types.NewChatID())
		require.Error(t, err)
		assert.Empty(t, p)

		p, err = jobresolveproblem.MarshalPayload(
			types.NewUserID(),
			types.NewRequestID(),
			types.NewMessageID(),
			types.ChatIDNil)
		require.Error(t, err)
		assert.Empty(t, p)

	})

}
