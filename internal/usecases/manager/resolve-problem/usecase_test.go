package resolveproblem_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	problemsrepo "github.com/Pickausernaame/chat-service/internal/repositories/problems"
	jobresolveproblem "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/job-resolve-problem"
	"github.com/Pickausernaame/chat-service/internal/testingh"
	"github.com/Pickausernaame/chat-service/internal/types"
	resolveproblem "github.com/Pickausernaame/chat-service/internal/usecases/manager/resolve-problem"
	resolveproblemmocks "github.com/Pickausernaame/chat-service/internal/usecases/manager/resolve-problem/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl    *gomock.Controller
	outbox  *resolveproblemmocks.MockoutboxService
	msgRepo *resolveproblemmocks.MockmessagesRepository
	prbRepo *resolveproblemmocks.MockproblemsRepository
	txtor   *resolveproblemmocks.Mocktransactor

	chatID    types.ChatID
	managerID types.UserID
	reqID     types.RequestID

	p *problemsrepo.Problem
	m *messagesrepo.Message

	uCase resolveproblem.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.msgRepo = resolveproblemmocks.NewMockmessagesRepository(s.ctrl)
	s.prbRepo = resolveproblemmocks.NewMockproblemsRepository(s.ctrl)
	s.outbox = resolveproblemmocks.NewMockoutboxService(s.ctrl)
	s.txtor = resolveproblemmocks.NewMocktransactor(s.ctrl)

	s.chatID = types.NewChatID()
	s.managerID = types.NewUserID()
	s.reqID = types.NewRequestID()

	s.p = &problemsrepo.Problem{
		ID:        types.NewProblemID(),
		ChatID:    s.chatID,
		ManagerID: s.managerID,
		CreatedAt: time.Now(),
	}

	s.m = &messagesrepo.Message{
		ID:                  types.NewMessageID(),
		ChatID:              s.chatID,
		InitialRequestID:    s.reqID,
		Body:                "somebody once told me",
		IsVisibleForClient:  true,
		IsVisibleForManager: false,
		IsBlocked:           false,
		IsService:           true,
		CreatedAt:           time.Now(),
	}

	var err error
	s.uCase, err = resolveproblem.New(resolveproblem.NewOptions(s.outbox, s.prbRepo, s.msgRepo, s.txtor))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestRequestValidationError() {
	// Arrange.
	req := resolveproblem.Request{}

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, resolveproblem.ErrInvalidRequest)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestGetProblemSomeError() {
	// Arrange.
	req := resolveproblem.Request{
		ChatID:    s.chatID,
		ManagerID: s.managerID,
		RequestID: s.reqID,
	}

	s.prbRepo.EXPECT().GetProblemByChatAndManagerIDs(gomock.Any(), req.ChatID, req.ManagerID).
		Return(nil, errors.New("some error"))

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.NotErrorIs(err, resolveproblem.ErrInvalidRequest)
	s.NotErrorIs(err, resolveproblem.ErrProblemNotFound)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestGetProblemNotFoundError() {
	// Arrange.
	req := resolveproblem.Request{
		ChatID:    s.chatID,
		ManagerID: s.managerID,
		RequestID: s.reqID,
	}

	s.prbRepo.EXPECT().GetProblemByChatAndManagerIDs(gomock.Any(), req.ChatID, req.ManagerID).
		Return(nil, problemsrepo.ErrProblemNotFound)

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.NotErrorIs(err, resolveproblem.ErrInvalidRequest)
	s.ErrorIs(err, resolveproblem.ErrProblemNotFound)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestCreateResolvedMessageError() {
	// Arrange.
	req := resolveproblem.Request{
		ChatID:    s.chatID,
		ManagerID: s.managerID,
		RequestID: s.reqID,
	}

	s.prbRepo.EXPECT().GetProblemByChatAndManagerIDs(gomock.Any(), req.ChatID, req.ManagerID).
		Return(s.p, nil)

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})

	s.msgRepo.EXPECT().CreateProblemResolvedMessage(gomock.Any(), req.ChatID, s.p.ID, req.RequestID).Return(nil, errors.New("some error"))

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.NotErrorIs(err, resolveproblem.ErrInvalidRequest)
	s.NotErrorIs(err, resolveproblem.ErrProblemNotFound)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestResolveProblemSomeError() {
	// Arrange.
	req := resolveproblem.Request{
		ChatID:    s.chatID,
		ManagerID: s.managerID,
		RequestID: s.reqID,
	}

	s.prbRepo.EXPECT().GetProblemByChatAndManagerIDs(gomock.Any(), req.ChatID, req.ManagerID).
		Return(s.p, nil)

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})

	s.msgRepo.EXPECT().CreateProblemResolvedMessage(gomock.Any(), req.ChatID, s.p.ID, req.RequestID).Return(s.m, nil)

	s.prbRepo.EXPECT().ResolveProblem(gomock.Any(), s.p.ID, s.managerID).Return(errors.New("some error"))

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.NotErrorIs(err, resolveproblem.ErrInvalidRequest)
	s.NotErrorIs(err, resolveproblem.ErrProblemNotFound)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestResolveProblemNotFoundError() {
	// Arrange.
	req := resolveproblem.Request{
		ChatID:    s.chatID,
		ManagerID: s.managerID,
		RequestID: s.reqID,
	}

	s.prbRepo.EXPECT().GetProblemByChatAndManagerIDs(gomock.Any(), req.ChatID, req.ManagerID).
		Return(s.p, nil)

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})

	s.msgRepo.EXPECT().CreateProblemResolvedMessage(gomock.Any(), req.ChatID, s.p.ID, req.RequestID).Return(s.m, nil)

	s.prbRepo.EXPECT().ResolveProblem(gomock.Any(), s.p.ID, s.managerID).Return(problemsrepo.ErrProblemNotFound)

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.NotErrorIs(err, resolveproblem.ErrInvalidRequest)
	s.ErrorIs(err, resolveproblem.ErrProblemNotFound)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestPutJobError() {
	// Arrange.
	req := resolveproblem.Request{
		ChatID:    s.chatID,
		ManagerID: s.managerID,
		RequestID: s.reqID,
	}

	s.prbRepo.EXPECT().GetProblemByChatAndManagerIDs(gomock.Any(), req.ChatID, req.ManagerID).
		Return(s.p, nil)

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})

	s.msgRepo.EXPECT().CreateProblemResolvedMessage(gomock.Any(), req.ChatID, s.p.ID, req.RequestID).Return(s.m, nil)

	s.prbRepo.EXPECT().ResolveProblem(gomock.Any(), s.p.ID, s.managerID).Return(nil)

	payload, err := jobresolveproblem.MarshalPayload(s.managerID, s.reqID, s.m.ID, s.chatID)
	s.Require().NoError(err)

	s.outbox.EXPECT().Put(gomock.Any(), jobresolveproblem.Name, payload, gomock.Any()).Return(types.JobIDNil, errors.New("some error"))
	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.NotErrorIs(err, resolveproblem.ErrInvalidRequest)
	s.NotErrorIs(err, resolveproblem.ErrProblemNotFound)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestCommitError() {
	// Arrange.
	req := resolveproblem.Request{
		ChatID:    s.chatID,
		ManagerID: s.managerID,
		RequestID: s.reqID,
	}

	s.prbRepo.EXPECT().GetProblemByChatAndManagerIDs(gomock.Any(), req.ChatID, req.ManagerID).
		Return(s.p, nil)

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		}).Return(errors.New("some error"))

	s.msgRepo.EXPECT().CreateProblemResolvedMessage(gomock.Any(), req.ChatID, s.p.ID, req.RequestID).Return(s.m, nil)

	s.prbRepo.EXPECT().ResolveProblem(gomock.Any(), s.p.ID, s.managerID).Return(nil)

	payload, err := jobresolveproblem.MarshalPayload(s.managerID, s.reqID, s.m.ID, s.chatID)
	s.Require().NoError(err)

	s.outbox.EXPECT().Put(gomock.Any(), jobresolveproblem.Name, payload, gomock.Any()).Return(types.NewJobID(), nil)
	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.NotErrorIs(err, resolveproblem.ErrInvalidRequest)
	s.NotErrorIs(err, resolveproblem.ErrProblemNotFound)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestCommitSuccess() {
	// Arrange.
	req := resolveproblem.Request{
		ChatID:    s.chatID,
		ManagerID: s.managerID,
		RequestID: s.reqID,
	}

	s.prbRepo.EXPECT().GetProblemByChatAndManagerIDs(gomock.Any(), req.ChatID, req.ManagerID).
		Return(s.p, nil)

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		}).Return(nil)

	s.msgRepo.EXPECT().CreateProblemResolvedMessage(gomock.Any(), req.ChatID, s.p.ID, req.RequestID).Return(s.m, nil)

	s.prbRepo.EXPECT().ResolveProblem(gomock.Any(), s.p.ID, s.managerID).Return(nil)

	payload, err := jobresolveproblem.MarshalPayload(s.managerID, s.reqID, s.m.ID, s.chatID)
	s.Require().NoError(err)

	s.outbox.EXPECT().Put(gomock.Any(), jobresolveproblem.Name, payload, gomock.Any()).Return(types.NewJobID(), nil)
	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().NoError(err)
	s.Empty(resp)
}
