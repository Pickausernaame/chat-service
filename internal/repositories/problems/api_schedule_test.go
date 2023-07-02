//go:build integration

package problemsrepo_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	problemsrepo "github.com/Pickausernaame/chat-service/internal/repositories/problems"
	"github.com/Pickausernaame/chat-service/internal/testingh"
	"github.com/Pickausernaame/chat-service/internal/types"
)

type ProblemsSchedulerRepoSuite struct {
	testingh.DBSuite
	repo *problemsrepo.Repo
}

func TestProblemsSchedulerRepoSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ProblemsRepoSuite{DBSuite: testingh.NewDBSuite("TestProblemsSchedulerRepoSuite")})
}

func (s *ProblemsSchedulerRepoSuite) SetupSuite() {
	s.DBSuite.SetupSuite()

	var err error

	s.repo, err = problemsrepo.New(problemsrepo.NewOptions(s.Database))
	s.Require().NoError(err)
}

func (s *ProblemsSchedulerRepoSuite) Test_GetUnassignedProblems() {
	chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(types.NewUserID()).Save(s.Ctx)
	s.Require().NoError(err)

	problem1, err := s.Database.Problem(s.Ctx).Create().
		SetChatID(chat.ID).
		Save(s.Ctx)
	s.Require().NoError(err)
	s.NotEmpty(problem1)

	chat, err = s.Database.Chat(s.Ctx).Create().SetClientID(types.NewUserID()).Save(s.Ctx)
	s.Require().NoError(err)

	problem2, err := s.Database.Problem(s.Ctx).Create().
		SetChatID(chat.ID).
		Save(s.Ctx)
	s.Require().NoError(err)
	s.NotEmpty(problem2)

	chat, err = s.Database.Chat(s.Ctx).Create().SetClientID(types.NewUserID()).Save(s.Ctx)
	s.Require().NoError(err)

	// ignore resolved problems
	problem3, err := s.Database.Problem(s.Ctx).Create().
		SetChatID(chat.ID).
		SetResolveAt(time.Now()).Save(s.Ctx)
	s.Require().NoError(err)
	s.NotEmpty(problem3)

	chat, err = s.Database.Chat(s.Ctx).Create().SetClientID(types.NewUserID()).Save(s.Ctx)
	s.Require().NoError(err)
	// ignore assigned problems
	problem4, err := s.Database.Problem(s.Ctx).Create().
		SetChatID(chat.ID).
		SetManagerID(types.NewUserID()).
		Save(s.Ctx)
	s.Require().NoError(err)
	s.NotEmpty(problem4)

	problems, err := s.repo.GetUnassignedProblems(s.Ctx)
	s.Run("finding problem, that not assigned to manager", func() {
		s.Require().NoError(err)
		s.Equal(2, len(problems))
	})
}

func (s *ProblemsSchedulerRepoSuite) Test_AssignManager() {
	s.Run("assign manager to problem", func() {
		managerID := types.NewUserID()
		clientID := types.NewUserID()

		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		problem1, err := s.Database.Problem(s.Ctx).Create().
			SetChatID(chat.ID).
			Save(s.Ctx)
		s.Require().NoError(err)
		s.NotEmpty(problem1)

		err = s.repo.AssignManager(s.Ctx, problem1.ID, managerID)
		s.Require().NoError(err)

		res, err := s.Database.Problem(s.Ctx).Get(s.Ctx, problem1.ID)
		s.Require().NoError(err)
		s.Equal(managerID, res.ManagerID)
	})
}
