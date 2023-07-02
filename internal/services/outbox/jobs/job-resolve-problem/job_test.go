package jobresolveproblem_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	eventstream "github.com/Pickausernaame/chat-service/internal/services/event-stream"
	jobresolveproblem "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/job-resolve-problem"
	jobresolveproblemmocks "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/job-resolve-problem/mocks"
	"github.com/Pickausernaame/chat-service/internal/types"
)

func TestJob_Handle(t *testing.T) {
	// Arrange.
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	msgRepo := jobresolveproblemmocks.NewMockmessageRepository(ctrl)
	chatRepo := jobresolveproblemmocks.NewMockchatRepository(ctrl)

	managerLoad := jobresolveproblemmocks.NewMockmanagerLoad(ctrl)
	eventStream := jobresolveproblemmocks.NewMockeventStream(ctrl)

	job, err := jobresolveproblem.New(jobresolveproblem.NewOptions(msgRepo, chatRepo, managerLoad, eventStream))
	require.NoError(t, err)

	clientID := types.NewUserID()
	managerID := types.NewUserID()
	msgID := types.NewMessageID()
	chatID := types.NewChatID()
	reqID := types.NewRequestID()

	const body = "Hello!"

	msg := messagesrepo.Message{
		ID:                  msgID,
		ChatID:              chatID,
		Body:                body,
		InitialRequestID:    reqID,
		CreatedAt:           time.Now(),
		IsVisibleForClient:  true,
		IsVisibleForManager: false,
		IsBlocked:           false,
		IsService:           false,
	}
	msgRepo.EXPECT().GetMessageByID(gomock.Any(), msgID).Return(&msg, nil)
	managerLoad.EXPECT().CanManagerTakeProblem(ctx, managerID).Return(true, nil)
	chatRepo.EXPECT().ClientIDByID(gomock.Any(), msg.ChatID).Return(clientID, nil)

	closeChatEvent := &eventstream.ChatClosedEvent{
		EventType:           eventstream.EventTypeChatClosedEvent,
		RequestID:           msg.InitialRequestID,
		ChatID:              msg.ChatID,
		CanTakeMoreProblems: true,
	}

	eventStream.EXPECT().Publish(gomock.Any(), managerID, closeChatEvent).Times(1).Return(nil)

	newMsgEvent := &eventstream.NewMessageEvent{
		EventType:   eventstream.EventTypeNewMessageEvent,
		RequestID:   msg.InitialRequestID,
		ChatID:      msg.ChatID,
		MessageID:   msg.ID,
		CreatedAt:   msg.CreatedAt,
		MessageBody: msg.Body,
		IsService:   msg.IsService,
	}
	eventStream.EXPECT().Publish(gomock.Any(), clientID, newMsgEvent).Times(1).Return(nil)

	// Action & assert.
	payload, err := jobresolveproblem.MarshalPayload(managerID, reqID, msgID, chatID)
	require.NoError(t, err)

	err = job.Handle(ctx, payload)
	require.NoError(t, err)
}
