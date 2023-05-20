package canreceiveproblems_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/Pickausernaame/chat-service/internal/testingh"
	"github.com/Pickausernaame/chat-service/internal/types"
	canreceiveproblems "github.com/Pickausernaame/chat-service/internal/usecases/manager/can-receive-problems"
	canreceiveproblemsmocks "github.com/Pickausernaame/chat-service/internal/usecases/manager/can-receive-problems/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl      *gomock.Controller
	mLoadMock *canreceiveproblemsmocks.MockmanagerLoadService
	mPoolMock *canreceiveproblemsmocks.MockmanagerPool
	uCase     canreceiveproblems.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.ContextSuite.SetupTest()

	s.mLoadMock = canreceiveproblemsmocks.NewMockmanagerLoadService(s.ctrl)
	s.mPoolMock = canreceiveproblemsmocks.NewMockmanagerPool(s.ctrl)
	var err error
	s.uCase, err = canreceiveproblems.New(canreceiveproblems.NewOptions(s.mLoadMock, s.mPoolMock))
	s.Require().NoError(err)
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestRequestValidationError() {
	// Arrange.
	req := canreceiveproblems.Request{}

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, canreceiveproblems.ErrInvalidRequest)
	s.False(resp.Result)
}

func (s *UseCaseSuite) TestContainsError() {
	// Arrange.
	req := canreceiveproblems.Request{
		ID:        types.NewRequestID(),
		ManagerID: types.NewUserID(),
	}
	s.mPoolMock.EXPECT().Contains(gomock.Any(), req.ManagerID).Return(false, errors.New("some error"))

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.False(resp.Result)
}

func (s *UseCaseSuite) TestContainsReturnsTrue() {
	// Arrange.
	req := canreceiveproblems.Request{
		ID:        types.NewRequestID(),
		ManagerID: types.NewUserID(),
	}

	s.mPoolMock.EXPECT().Contains(gomock.Any(), req.ManagerID).Return(true, nil)

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().NoError(err)
	s.False(resp.Result)
}

func (s *UseCaseSuite) TestCanManagerTakeProblemError() {
	// Arrange.
	req := canreceiveproblems.Request{
		ID:        types.NewRequestID(),
		ManagerID: types.NewUserID(),
	}

	s.mPoolMock.EXPECT().Contains(gomock.Any(), req.ManagerID).Return(false, nil)

	s.mLoadMock.EXPECT().CanManagerTakeProblem(gomock.Any(), req.ManagerID).
		Return(false, errors.New("some error"))

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.False(resp.Result)
}

func (s *UseCaseSuite) TestSuccess() {
	// Arrange.
	req := canreceiveproblems.Request{
		ID:        types.NewRequestID(),
		ManagerID: types.NewUserID(),
	}

	s.mPoolMock.EXPECT().Contains(gomock.Any(), req.ManagerID).Return(false, nil)

	s.mLoadMock.EXPECT().CanManagerTakeProblem(gomock.Any(), req.ManagerID).
		Return(false, nil)

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().NoError(err)
	s.False(resp.Result)

	// Arrange.
	s.mPoolMock.EXPECT().Contains(gomock.Any(), req.ManagerID).Return(false, nil)

	s.mLoadMock.EXPECT().CanManagerTakeProblem(gomock.Any(), req.ManagerID).
		Return(true, nil)

	// Action.
	resp, err = s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().NoError(err)
	s.True(resp.Result)
}
