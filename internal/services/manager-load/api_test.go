package managerload_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	managerload "github.com/Pickausernaame/chat-service/internal/services/manager-load"
	managerloadmocks "github.com/Pickausernaame/chat-service/internal/services/manager-load/mocks"
	"github.com/Pickausernaame/chat-service/internal/testingh"
	"github.com/Pickausernaame/chat-service/internal/types"
)

type ServiceSuite struct {
	testingh.ContextSuite

	ctrl *gomock.Controller

	problemsRepo *managerloadmocks.MockproblemsRepository
	managerLoad  *managerload.Service
}

func TestServiceSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ServiceSuite))
}

func (s *ServiceSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.problemsRepo = managerloadmocks.NewMockproblemsRepository(s.ctrl)

	// check options validation
	var err error
	cases := []struct {
		wantErr     bool
		description string
		maxProblems int
	}{
		{
			true,
			"negative problems count",
			-1,
		},
		{
			true,
			"max problems limit exceeded",
			31,
		},
		{
			false,
			"min problems count",
			1,
		},
		{
			false,
			"max problems count",
			30,
		},
	}
	for _, c := range cases {
		_, err := managerload.New(managerload.NewOptions(c.maxProblems, s.problemsRepo))
		if c.wantErr {
			s.Require().Error(err)
		} else {
			s.Require().NoError(err)
		}
	}

	s.managerLoad, err = managerload.New(managerload.NewOptions(2, s.problemsRepo))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *ServiceSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *ServiceSuite) Test_CanManagerTakeProblem() {
	ctx := context.Background()
	managerID := types.NewUserID()

	// success
	s.problemsRepo.EXPECT().GetManagerOpenProblemsCount(gomock.Any(), managerID).Return(1, nil)
	isCan, err := s.managerLoad.CanManagerTakeProblem(ctx, managerID)
	s.Require().NoError(err)
	s.Require().True(isCan)

	// negative
	s.problemsRepo.EXPECT().GetManagerOpenProblemsCount(gomock.Any(), managerID).Return(10, nil)
	isCan, err = s.managerLoad.CanManagerTakeProblem(ctx, managerID)
	s.Require().NoError(err)
	s.Require().False(isCan)

	// with error
	s.problemsRepo.EXPECT().GetManagerOpenProblemsCount(gomock.Any(), managerID).Return(0, errors.New("some error"))
	isCan, err = s.managerLoad.CanManagerTakeProblem(ctx, managerID)
	s.Require().Error(err)
	s.Require().False(isCan)
}
