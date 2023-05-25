package sendclientmessagejob_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	eventstream "github.com/Pickausernaame/chat-service/internal/services/event-stream"
	msgproducer "github.com/Pickausernaame/chat-service/internal/services/msg-producer"
	sendclientmessagejob "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/send-client-message"
	sendclientmessagejobmocks "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/send-client-message/mocks"
	"github.com/Pickausernaame/chat-service/internal/types"
)

func TestJob_Handle(t *testing.T) {
	// Arrange.
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	msgProducer := sendclientmessagejobmocks.NewMockmessageProducer(ctrl)
	msgRepo := sendclientmessagejobmocks.NewMockmessageRepository(ctrl)
	eventStream := sendclientmessagejobmocks.NewMockeventStream(ctrl)
	job, err := sendclientmessagejob.New(sendclientmessagejob.NewOptions(msgProducer, msgRepo, eventStream))
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
		IsService:           false,
	}
	msgRepo.EXPECT().GetMessageByID(gomock.Any(), msgID).Return(&msg, nil)

	msgProducer.EXPECT().ProduceMessage(gomock.Any(), msgproducer.Message{
		ID:         msgID,
		ChatID:     chatID,
		Body:       body,
		FromClient: true,
	}).Return(nil)

	event := &eventstream.NewMessageEvent{
		EventType:   eventstream.EventTypeNewMessageEvent,
		RequestID:   msg.InitialRequestID,
		ChatID:      msg.ChatID,
		MessageID:   msg.ID,
		UserID:      msg.AuthorID,
		CreatedAt:   msg.CreatedAt,
		MessageBody: msg.Body,
		IsService:   msg.IsService,
	}

	eventStream.EXPECT().Publish(gomock.Any(), msg.AuthorID, event)
	// Action & assert.
	payload, err := sendclientmessagejob.MarshalPayload(msgID)
	require.NoError(t, err)

	err = job.Handle(ctx, payload)
	require.NoError(t, err)
}
