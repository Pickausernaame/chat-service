package clientmessageblockedjob_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	eventstream "github.com/Pickausernaame/chat-service/internal/services/event-stream"
	clientmessageblockedjob "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/client-message-blocked"
	clientmessageblockedjobmocks "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/client-message-blocked/mocks"
	"github.com/Pickausernaame/chat-service/internal/types"
)

func TestJob_Handle(t *testing.T) {
	// Arrange.
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	msgRepo := clientmessageblockedjobmocks.NewMockmessageRepository(ctrl)
	eventStream := clientmessageblockedjobmocks.NewMockeventStream(ctrl)
	job, err := clientmessageblockedjob.New(clientmessageblockedjob.NewOptions(msgRepo, eventStream))
	require.NoError(t, err)

	clientID := types.NewUserID()
	msgID := types.NewMessageID()
	chatID := types.NewChatID()
	const body = "Hello!"

	msg := messagesrepo.Message{
		ID:                  msgID,
		ChatID:              chatID,
		AuthorID:            clientID,
		Body:                body,
		CreatedAt:           time.Now(),
		IsVisibleForClient:  true,
		IsVisibleForManager: false,
		IsBlocked:           false,
		InitialRequestID:    types.NewRequestID(),
		IsService:           false,
	}
	msgRepo.EXPECT().GetMessageByID(gomock.Any(), msgID).Return(&msg, nil)

	event := &eventstream.MessageBlockedEvent{
		EventType: eventstream.EventTypeMessageBlockedEvent,
		RequestID: msg.InitialRequestID,
		MessageID: msg.ID,
	}

	eventStream.EXPECT().Publish(gomock.Any(), msg.AuthorID, event)
	// Action & assert.
	payload, err := clientmessageblockedjob.MarshalPayload(msgID)
	require.NoError(t, err)

	err = job.Handle(ctx, payload)
	require.NoError(t, err)
}
