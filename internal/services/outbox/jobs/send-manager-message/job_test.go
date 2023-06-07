package sendmanagermessagejob_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	eventstream "github.com/Pickausernaame/chat-service/internal/services/event-stream"
	msgproducer "github.com/Pickausernaame/chat-service/internal/services/msg-producer"
	sendmanagermessagejob "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/send-manager-message"
	sendmanagermessagejobmocks "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/send-manager-message/mocks"
	"github.com/Pickausernaame/chat-service/internal/types"
)

func TestJob_Handle(t *testing.T) {
	// Arrange.
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	msgProducer := sendmanagermessagejobmocks.NewMockmessageProducer(ctrl)
	msgRepo := sendmanagermessagejobmocks.NewMockmessageRepository(ctrl)
	chatRepo := sendmanagermessagejobmocks.NewMockchatRepository(ctrl)
	eventStream := sendmanagermessagejobmocks.NewMockeventStream(ctrl)
	job, err := sendmanagermessagejob.New(sendmanagermessagejob.NewOptions(msgProducer, msgRepo, chatRepo, eventStream))
	require.NoError(t, err)

	clientID := types.NewUserID()
	managerID := types.NewUserID()
	msgID := types.NewMessageID()
	chatID := types.NewChatID()
	const body = "Hello!"

	msg := messagesrepo.Message{
		ID:                  msgID,
		ChatID:              chatID,
		AuthorID:            managerID,
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
		FromClient: false,
	}).Return(nil)

	chatRepo.EXPECT().ClientIDByID(gomock.Any(), msg.ChatID).Return(clientID, nil)

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

	eventStream.EXPECT().Publish(gomock.Any(), clientID, event)
	// Action & assert.
	payload, err := sendmanagermessagejob.MarshalPayload(msgID)
	require.NoError(t, err)

	err = job.Handle(ctx, payload)
	require.NoError(t, err)
}
