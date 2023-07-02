package getassignedproblems_test

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	problemsrepo "github.com/Pickausernaame/chat-service/internal/repositories/problems"
	"github.com/Pickausernaame/chat-service/internal/testingh"
	"github.com/Pickausernaame/chat-service/internal/types"
	getassignedproblems "github.com/Pickausernaame/chat-service/internal/usecases/manager/get-assigned-problems"
	getassignedproblemsmocks "github.com/Pickausernaame/chat-service/internal/usecases/manager/get-assigned-problems/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl            *gomock.Controller
	chatRepoMock    *getassignedproblemsmocks.MockchatRepository
	problemRepoMock *getassignedproblemsmocks.MockproblemRepository
	uCase           getassignedproblems.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.ContextSuite.SetupTest()

	s.chatRepoMock = getassignedproblemsmocks.NewMockchatRepository(s.ctrl)
	s.problemRepoMock = getassignedproblemsmocks.NewMockproblemRepository(s.ctrl)
	var err error
	s.uCase, err = getassignedproblems.New(getassignedproblems.NewOptions(s.problemRepoMock, s.chatRepoMock))
	s.Require().NoError(err)
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestRequestValidationError() {
	// Arrange.
	req := getassignedproblems.Request{}

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, getassignedproblems.ErrInvalidRequest)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestGetAssignedUnsolvedProblemsError() {
	// Arrange.
	req := getassignedproblems.Request{
		ManagerID: types.NewUserID(),
	}
	s.problemRepoMock.EXPECT().GetAssignedUnsolvedProblems(gomock.Any(), req.ManagerID).Return(nil, errors.New("some error"))

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestHandleSuccess() {
	// Arrange.
	req := getassignedproblems.Request{
		ManagerID: types.NewUserID(),
	}

	problems := make([]*problemsrepo.ProblemAndClientID, 0, 10)
	expectedResponse := make([]*getassignedproblems.Chat, 0, 10)
	for i := 0; i < 10; i++ {
		p := &problemsrepo.ProblemAndClientID{
			Problem: &problemsrepo.Problem{
				ID:        types.NewProblemID(),
				ChatID:    types.NewChatID(),
				ManagerID: req.ManagerID,
				CreatedAt: time.Now(),
			},
			ClientID: types.NewUserID(),
		}
		problems = append(problems, p)

		expectedResponse = append(expectedResponse, &getassignedproblems.Chat{
			ChatID:   p.ChatID,
			ClientID: p.ClientID,
		})
	}

	s.problemRepoMock.EXPECT().GetAssignedUnsolvedProblems(gomock.Any(), req.ManagerID).Return(problems, nil)

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().NoError(err)
	s.Equal(expectedResponse, resp.Chats)
}
