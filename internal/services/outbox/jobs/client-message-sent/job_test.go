package clientmessagesentjob_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	eventstream "github.com/Pickausernaame/chat-service/internal/services/event-stream"
	clientmessagesentjob "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/client-message-sent"
	clientmessagesentjobmocks "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/client-message-sent/mocks"
	"github.com/Pickausernaame/chat-service/internal/testingh"
	"github.com/Pickausernaame/chat-service/internal/types"
)

type JobSuite struct {
	testingh.ContextSuite

	ctrl        *gomock.Controller
	msgRepo     *clientmessagesentjobmocks.MockmessageRepository
	problemRepo *clientmessagesentjobmocks.MockproblemRepository
	eventStream *clientmessagesentjobmocks.MockeventStream

	job *clientmessagesentjob.Job
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(JobSuite))
}

func (s *JobSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.msgRepo = clientmessagesentjobmocks.NewMockmessageRepository(s.ctrl)
	s.problemRepo = clientmessagesentjobmocks.NewMockproblemRepository(s.ctrl)
	s.eventStream = clientmessagesentjobmocks.NewMockeventStream(s.ctrl)

	var err error
	s.job, err = clientmessagesentjob.New(clientmessagesentjob.NewOptions(s.msgRepo, s.problemRepo, s.eventStream))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *JobSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *JobSuite) Test_HandleSuccessWithManagerEvent() {
	// Arrange.
	clientID := types.NewUserID()
	managerID := types.NewUserID()
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
	s.msgRepo.EXPECT().GetMessageByID(gomock.Any(), msgID).Return(&msg, nil)
	// expect sending event to client
	sentEvent := &eventstream.MessageSentEvent{
		EventType: eventstream.EventTypeMessageSentEvent,
		RequestID: msg.InitialRequestID,
		MessageID: msg.ID,
	}
	s.eventStream.EXPECT().Publish(gomock.Any(), msg.AuthorID, sentEvent)

	s.problemRepo.EXPECT().GetManagerIDByChatID(gomock.Any(), chatID).Return(managerID, nil)

	newMessageEvent := &eventstream.NewMessageEvent{
		EventType:   eventstream.EventTypeNewMessageEvent,
		RequestID:   msg.InitialRequestID,
		ChatID:      msg.ChatID,
		MessageID:   msg.ID,
		UserID:      msg.AuthorID,
		CreatedAt:   msg.CreatedAt,
		MessageBody: msg.Body,
		IsService:   msg.IsService,
	}
	s.eventStream.EXPECT().Publish(gomock.Any(), managerID, newMessageEvent)

	// Action & assert.
	payload, err := clientmessagesentjob.MarshalPayload(msgID)
	s.Require().NoError(err)

	err = s.job.Handle(s.Ctx, payload)
	s.Require().NoError(err)
}

func (s *JobSuite) Test_HandleSuccessWithoutManagerEvent() {
	// Arrange.
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
	s.msgRepo.EXPECT().GetMessageByID(gomock.Any(), msgID).Return(&msg, nil)
	// expect sending event to client
	sentEvent := &eventstream.MessageSentEvent{
		EventType: eventstream.EventTypeMessageSentEvent,
		RequestID: msg.InitialRequestID,
		MessageID: msg.ID,
	}
	s.eventStream.EXPECT().Publish(gomock.Any(), msg.AuthorID, sentEvent)

	s.problemRepo.EXPECT().GetManagerIDByChatID(gomock.Any(), chatID).Return(types.UserIDNil, nil)

	// Action & assert.
	payload, err := clientmessagesentjob.MarshalPayload(msgID)
	s.Require().NoError(err)

	err = s.job.Handle(s.Ctx, payload)
	s.Require().NoError(err)
}
