package sendmessage_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	problemsrepo "github.com/Pickausernaame/chat-service/internal/repositories/problems"
	sendmanagermessagejob "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/send-manager-message"
	"github.com/Pickausernaame/chat-service/internal/testingh"
	"github.com/Pickausernaame/chat-service/internal/types"
	sendmessage "github.com/Pickausernaame/chat-service/internal/usecases/manager/send-message"
	sendmessagemocks "github.com/Pickausernaame/chat-service/internal/usecases/manager/send-message/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl        *gomock.Controller
	msgRepo     *sendmessagemocks.MockmessagesRepository
	problemRepo *sendmessagemocks.MockproblemsRepository
	txtor       *sendmessagemocks.Mocktransactor
	outBoxSvc   *sendmessagemocks.MockoutboxService
	uCase       sendmessage.UseCase

	msg *messagesrepo.Message
	req sendmessage.Request
	prb *problemsrepo.Problem
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.msgRepo = sendmessagemocks.NewMockmessagesRepository(s.ctrl)
	s.outBoxSvc = sendmessagemocks.NewMockoutboxService(s.ctrl)
	s.problemRepo = sendmessagemocks.NewMockproblemsRepository(s.ctrl)
	s.txtor = sendmessagemocks.NewMocktransactor(s.ctrl)

	s.req = sendmessage.Request{
		ID:          types.NewRequestID(),
		ManagerID:   types.NewUserID(),
		ChatID:      types.NewChatID(),
		MessageBody: "где деньги?",
	}

	s.msg = &messagesrepo.Message{
		ID:                  types.NewMessageID(),
		ChatID:              s.req.ChatID,
		AuthorID:            s.req.ManagerID,
		InitialRequestID:    s.req.ID,
		Body:                s.req.MessageBody,
		IsVisibleForClient:  true,
		IsVisibleForManager: true,
		CreatedAt:           time.Now(),
	}

	s.prb = &problemsrepo.Problem{
		ID:        types.NewProblemID(),
		ChatID:    s.req.ChatID,
		ManagerID: s.req.ManagerID,
		CreatedAt: time.Now(),
	}

	var err error
	s.uCase, err = sendmessage.New(sendmessage.NewOptions(s.msgRepo, s.outBoxSvc, s.problemRepo, s.txtor))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestRequestValidationError() {
	// Arrange.
	req := sendmessage.Request{}

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, sendmessage.ErrInvalidRequest)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestGetProblemError() {
	// Arrange.

	s.problemRepo.EXPECT().GetProblemByChatAndManagerIDs(gomock.Any(), s.req.ChatID, s.req.ManagerID).
		Return(nil, errors.New("some error"))

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, s.req)

	// Assert.
	s.Require().Error(err)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestCreateMessageError() {
	// Arrange.

	s.problemRepo.EXPECT().GetProblemByChatAndManagerIDs(gomock.Any(), s.req.ChatID, s.req.ManagerID).
		Return(s.prb, nil)

	err := errors.New("some error")
	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})

	s.msgRepo.EXPECT().CreateFullVisible(gomock.Any(), s.req.ID, s.prb.ID, s.prb.ChatID, s.prb.ManagerID, s.req.MessageBody).
		Return(nil, err)
	// Action.
	resp, err := s.uCase.Handle(s.Ctx, s.req)

	// Assert.
	s.Require().Error(err)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestCreateJobError() {
	// Arrange.

	s.problemRepo.EXPECT().GetProblemByChatAndManagerIDs(gomock.Any(), s.req.ChatID, s.req.ManagerID).
		Return(s.prb, nil)

	resErr := errors.New("some error")
	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})

	s.msgRepo.EXPECT().CreateFullVisible(gomock.Any(), s.req.ID, s.prb.ID, s.prb.ChatID, s.prb.ManagerID, s.req.MessageBody).
		Return(s.msg, nil)

	payload, err := sendmanagermessagejob.MarshalPayload(s.msg.ID)
	s.Require().NoError(err)

	s.outBoxSvc.EXPECT().Put(gomock.Any(), sendmanagermessagejob.Name, payload, gomock.Any()).
		Return(types.JobIDNil, resErr)
	// Action.
	resp, err := s.uCase.Handle(s.Ctx, s.req)

	// Assert.
	s.Require().Error(err)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestCommitError() {
	// Arrange.

	s.problemRepo.EXPECT().GetProblemByChatAndManagerIDs(gomock.Any(), s.req.ChatID, s.req.ManagerID).
		Return(s.prb, nil)

	resErr := errors.New("some error")
	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		}).Return(resErr)

	s.msgRepo.EXPECT().CreateFullVisible(gomock.Any(), s.req.ID, s.prb.ID, s.prb.ChatID, s.prb.ManagerID, s.req.MessageBody).
		Return(s.msg, nil)

	payload, err := sendmanagermessagejob.MarshalPayload(s.msg.ID)
	s.Require().NoError(err)

	s.outBoxSvc.EXPECT().Put(gomock.Any(), sendmanagermessagejob.Name, payload, gomock.Any()).
		Return(types.NewJobID(), nil)
	// Action.
	resp, err := s.uCase.Handle(s.Ctx, s.req)

	// Assert.
	s.Require().Error(err)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestCommitSuccess() {
	// Arrange.

	s.problemRepo.EXPECT().GetProblemByChatAndManagerIDs(gomock.Any(), s.req.ChatID, s.req.ManagerID).
		Return(s.prb, nil)

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})

	s.msgRepo.EXPECT().CreateFullVisible(gomock.Any(), s.req.ID, s.prb.ID, s.prb.ChatID, s.prb.ManagerID, s.req.MessageBody).
		Return(s.msg, nil)

	payload, err := sendmanagermessagejob.MarshalPayload(s.msg.ID)
	s.Require().NoError(err)

	s.outBoxSvc.EXPECT().Put(gomock.Any(), sendmanagermessagejob.Name, payload, gomock.Any()).
		Return(types.NewJobID(), nil)
	// Action.
	resp, err := s.uCase.Handle(s.Ctx, s.req)

	// Assert.
	s.Require().NoError(err)
	s.Equal(s.msg.ID, resp.MessageID)
	s.Equal(s.msg.CreatedAt, resp.CreatedAt)
}
