package managerassignedtoproblemjob_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	eventstream "github.com/Pickausernaame/chat-service/internal/services/event-stream"
	managerassignedtoproblemjob "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/manager-assigned-to-problem"
	managerassignedtoproblemjobmocks "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/manager-assigned-to-problem/mocks"
	"github.com/Pickausernaame/chat-service/internal/testingh"
	"github.com/Pickausernaame/chat-service/internal/types"
)

type JobSuite struct {
	testingh.ContextSuite

	ctrl        *gomock.Controller
	msgRepo     *managerassignedtoproblemjobmocks.MockmessageRepository
	managerLoad *managerassignedtoproblemjobmocks.MockmanagerLoad
	eventStream *managerassignedtoproblemjobmocks.MockeventStream

	job *managerassignedtoproblemjob.Job
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(JobSuite))
}

func (s *JobSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.msgRepo = managerassignedtoproblemjobmocks.NewMockmessageRepository(s.ctrl)
	s.managerLoad = managerassignedtoproblemjobmocks.NewMockmanagerLoad(s.ctrl)
	s.eventStream = managerassignedtoproblemjobmocks.NewMockeventStream(s.ctrl)

	var err error
	s.job, err = managerassignedtoproblemjob.New(
		managerassignedtoproblemjob.NewOptions(s.msgRepo, s.managerLoad, s.eventStream))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *JobSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *JobSuite) Test_HandleSuccess() {
	// Arrange.
	clientID := types.NewUserID()
	managerID := types.NewUserID()
	reqID := types.NewRequestID()
	msgID := types.NewMessageID()
	chatID := types.NewChatID()
	const isCanTake = true
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
		InitialRequestID:    reqID,
		IsService:           false,
	}

	s.msgRepo.EXPECT().GetMessageByID(gomock.Any(), msgID).Return(&msg, nil)
	s.managerLoad.EXPECT().CanManagerTakeProblem(s.Ctx, managerID).Return(isCanTake, nil)
	// expect sending event to client
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
	s.eventStream.EXPECT().Publish(gomock.Any(), msg.AuthorID, newMessageEvent)

	// expect sending event to manager
	newChatEvent := &eventstream.NewChatEvent{
		EventType:           eventstream.EventTypeNewChatEvent,
		RequestID:           msg.InitialRequestID,
		ChatID:              msg.ChatID,
		CanTakeMoreProblems: isCanTake,
		ClientID:            clientID,
	}
	s.eventStream.EXPECT().Publish(gomock.Any(), managerID, newChatEvent)

	// Action & assert.
	payload, err := managerassignedtoproblemjob.MarshalPayload(clientID, managerID, reqID, msgID)
	s.Require().NoError(err)

	err = s.job.Handle(s.Ctx, payload)
	s.Require().NoError(err)
}
