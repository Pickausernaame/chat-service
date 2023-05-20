package setreadyreceiveproblems_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/Pickausernaame/chat-service/internal/testingh"
	"github.com/Pickausernaame/chat-service/internal/types"
	setreadyreceiveproblems "github.com/Pickausernaame/chat-service/internal/usecases/manager/set-ready-receive-problems"
	setreadyreceiveproblemsmocks "github.com/Pickausernaame/chat-service/internal/usecases/manager/set-ready-receive-problems/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl      *gomock.Controller
	mLoadMock *setreadyreceiveproblemsmocks.MockmanagerLoadService
	mPoolMock *setreadyreceiveproblemsmocks.MockmanagerPool
	uCase     setreadyreceiveproblems.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.ContextSuite.SetupTest()

	s.mLoadMock = setreadyreceiveproblemsmocks.NewMockmanagerLoadService(s.ctrl)
	s.mPoolMock = setreadyreceiveproblemsmocks.NewMockmanagerPool(s.ctrl)
	var err error
	s.uCase, err = setreadyreceiveproblems.New(setreadyreceiveproblems.NewOptions(s.mLoadMock, s.mPoolMock))
	s.Require().NoError(err)
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestRequestValidationError() {
	// Arrange.
	req := setreadyreceiveproblems.Request{}

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, setreadyreceiveproblems.ErrInvalidRequest)
}

func (s *UseCaseSuite) TestOverloadError() {
	// Arrange.
	req := setreadyreceiveproblems.Request{
		ID:        types.NewRequestID(),
		ManagerID: types.NewUserID(),
	}
	s.mLoadMock.EXPECT().CanManagerTakeProblem(gomock.Any(), req.ManagerID).
		Return(false, nil)
	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, setreadyreceiveproblems.ErrManagerOverload)
}

func (s *UseCaseSuite) TestCanManagerTakeProblemError() {
	// Arrange.
	req := setreadyreceiveproblems.Request{
		ID:        types.NewRequestID(),
		ManagerID: types.NewUserID(),
	}

	s.mLoadMock.EXPECT().CanManagerTakeProblem(gomock.Any(), req.ManagerID).
		Return(false, errors.New("some error"))

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestPutManagerToPoolError() {
	// Arrange.
	req := setreadyreceiveproblems.Request{
		ID:        types.NewRequestID(),
		ManagerID: types.NewUserID(),
	}

	s.mLoadMock.EXPECT().CanManagerTakeProblem(gomock.Any(), req.ManagerID).
		Return(true, nil)

	s.mPoolMock.EXPECT().Put(gomock.Any(), req.ManagerID).Return(errors.New("some error"))

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestSuccess() {
	// Arrange.
	req := setreadyreceiveproblems.Request{
		ID:        types.NewRequestID(),
		ManagerID: types.NewUserID(),
	}

	s.mLoadMock.EXPECT().CanManagerTakeProblem(gomock.Any(), req.ManagerID).
		Return(true, nil)

	s.mPoolMock.EXPECT().Put(gomock.Any(), req.ManagerID).Return(nil)

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().NoError(err)
}
