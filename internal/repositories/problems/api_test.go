//go:build integration

package problemsrepo_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	problemsrepo "github.com/Pickausernaame/chat-service/internal/repositories/problems"
	storeproblem "github.com/Pickausernaame/chat-service/internal/store/problem"
	"github.com/Pickausernaame/chat-service/internal/testingh"
	"github.com/Pickausernaame/chat-service/internal/types"
)

type ProblemsRepoSuite struct {
	testingh.DBSuite
	repo *problemsrepo.Repo
}

func TestProblemsRepoSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ProblemsRepoSuite{DBSuite: testingh.NewDBSuite("TestProblemsRepoSuite")})
}

func (s *ProblemsRepoSuite) SetupSuite() {
	s.DBSuite.SetupSuite()

	var err error

	s.repo, err = problemsrepo.New(problemsrepo.NewOptions(s.Database))
	s.Require().NoError(err)
}

func (s *ProblemsRepoSuite) Test_CreateIfNotExists() {
	s.Run("problem does not exist, should be created", func() {
		clientID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		problemID, err := s.repo.CreateIfNotExists(s.Ctx, chat.ID)
		s.Require().NoError(err)
		s.NotEmpty(problemID)

		problem, err := s.Database.Problem(s.Ctx).Get(s.Ctx, problemID)
		s.Require().NoError(err)
		s.Equal(problemID, problem.ID)
		s.Equal(chat.ID, problem.ChatID)
	})

	s.Run("resolved problem already exists, should be created", func() {
		clientID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		// Create problem.
		problem, err := s.Database.Problem(s.Ctx).Create().
			SetChatID(chat.ID).
			SetManagerID(types.NewUserID()).
			SetResolveAt(time.Now()).Save(s.Ctx)
		s.Require().NoError(err)

		problemID, err := s.repo.CreateIfNotExists(s.Ctx, chat.ID)
		s.Require().NoError(err)
		s.NotEmpty(problemID)
		s.NotEqual(problem.ID, problemID)
	})

	s.Run("problem already exists", func() {
		clientID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		// Create problem.
		problem, err := s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).Save(s.Ctx)
		s.Require().NoError(err)

		problemID, err := s.repo.CreateIfNotExists(s.Ctx, chat.ID)
		s.Require().NoError(err)
		s.NotEmpty(problemID)
		s.Equal(problem.ID, problemID)
	})
}

func (s *ProblemsRepoSuite) Test_GetManagerOpenProblemsCount() {
	s.Run("manager has no open problems", func() {
		managerID := types.NewUserID()

		count, err := s.repo.GetManagerOpenProblemsCount(s.Ctx, managerID)
		s.Require().NoError(err)
		s.Empty(count)
	})

	s.Run("manager has open problems", func() {
		const (
			problemsCount         = 20
			resolvedProblemsCount = 3
		)

		managerID := types.NewUserID()
		problems := make([]types.ProblemID, 0, problemsCount)

		for i := 0; i < problemsCount; i++ {
			_, pID := s.createChatWithProblemAssignedTo(managerID)
			problems = append(problems, pID)
		}

		// Create problems for other managers.
		for i := 0; i < problemsCount; i++ {
			s.createChatWithProblemAssignedTo(types.NewUserID())
		}

		count, err := s.repo.GetManagerOpenProblemsCount(s.Ctx, managerID)
		s.Require().NoError(err)
		s.Equal(problemsCount, count)

		// Resolve some problems.
		for i := 0; i < resolvedProblemsCount; i++ {
			pID := problems[i*resolvedProblemsCount]
			_, err := s.Database.Problem(s.Ctx).
				Update().
				Where(storeproblem.ID(pID)).
				SetResolveAt(time.Now()).
				Save(s.Ctx)
			s.Require().NoError(err)
		}

		count, err = s.repo.GetManagerOpenProblemsCount(s.Ctx, managerID)
		s.Require().NoError(err)
		s.Equal(problemsCount-resolvedProblemsCount, count)
	})
}

func (s *ProblemsRepoSuite) Test_GetAssignedUnsolvedProblems() {
	s.Run("getting assigned problems to manager", func() {
		managerID := types.NewUserID()
		problems := make([]*problemsrepo.Problem, 0, 10)

		for i := 0; i < 10; i++ {
			clientID := types.NewUserID()
			chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
			s.Require().NoError(err)

			problem, err := s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).SetManagerID(managerID).Save(s.Ctx)
			s.Require().NoError(err)

			p := &problemsrepo.Problem{
				ID:        problem.ID,
				ChatID:    problem.ChatID,
				ManagerID: problem.ManagerID,
				ResolveAt: problem.ResolveAt,
				CreatedAt: problem.CreatedAt,
			}
			problems = append(problems, p)
		}

		// add resolved problems - should be ignored
		for i := 0; i < 3; i++ {
			clientID := types.NewUserID()
			chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
			s.Require().NoError(err)

			_, err = s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).SetManagerID(managerID).
				SetResolveAt(time.Now()).Save(s.Ctx)
			s.Require().NoError(err)
		}

		// add unsolved problems another manager - should be ignored
		anotherManagerID := types.NewUserID()
		for i := 0; i < 3; i++ {
			clientID := types.NewUserID()
			chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
			s.Require().NoError(err)

			_, err = s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).SetManagerID(anotherManagerID).Save(s.Ctx)
			s.Require().NoError(err)
		}

		// add solved problems another manager - should be ignored
		for i := 0; i < 3; i++ {
			clientID := types.NewUserID()
			chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
			s.Require().NoError(err)

			_, err = s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).SetManagerID(anotherManagerID).
				SetResolveAt(time.Now()).Save(s.Ctx)
			s.Require().NoError(err)
		}

		res, err := s.repo.GetAssignedUnsolvedProblems(s.Ctx, managerID)
		s.Require().NoError(err)
		s.Equal(len(problems), len(res))
		for i := 0; i < 10; i++ {
			s.Equal(problems[i].ChatID, res[i].ChatID)
			s.Equal(problems[i].ID, res[i].ID)
			s.Equal(problems[i].ManagerID, res[i].ManagerID)
		}
	})
}

func (s *ProblemsRepoSuite) Test_GetManagerIDByChatID() {
	s.Run("getting manager id by chat id - success", func() {
		managerID := types.NewUserID()

		clientID := types.NewUserID()
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		_, err = s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).SetManagerID(managerID).Save(s.Ctx)
		s.Require().NoError(err)

		res, err := s.repo.GetManagerIDByChatID(s.Ctx, chat.ID)
		s.Require().NoError(err)
		s.Equal(managerID, res)
	})

	s.Run("chat not exist", func() {
		managerID := types.NewUserID()
		chatID := types.NewChatID()
		clientID := types.NewUserID()
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		_, err = s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).SetManagerID(managerID).Save(s.Ctx)
		s.Require().NoError(err)

		res, err := s.repo.GetManagerIDByChatID(s.Ctx, chatID)
		s.Require().Error(err)
		s.Empty(res)
	})
}

func (s *ProblemsRepoSuite) createChatWithProblemAssignedTo(managerID types.UserID) (types.ChatID, types.ProblemID) {
	s.T().Helper()

	// 1 chat can have only 1 open problem.

	chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(types.NewUserID()).Save(s.Ctx)
	s.Require().NoError(err)

	p, err := s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).SetManagerID(managerID).Save(s.Ctx)
	s.Require().NoError(err)

	return chat.ID, p.ID
}

func (s *ProblemsRepoSuite) Test_ResolveProblem() {
	s.Run("resolve problem", func() {
		clientID := types.NewUserID()
		managerID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		problem1, err := s.Database.Problem(s.Ctx).Create().
			SetChatID(chat.ID).
			SetManagerID(managerID).
			Save(s.Ctx)
		s.Require().NoError(err)
		s.NotEmpty(problem1)

		res, err := s.Database.Problem(s.Ctx).Get(s.Ctx, problem1.ID)
		s.Require().NoError(err)
		s.Empty(res.ResolveAt)

		err = s.repo.ResolveProblem(s.Ctx, problem1.ID, managerID)
		s.Require().NoError(err)

		res, err = s.Database.Problem(s.Ctx).Get(s.Ctx, problem1.ID)
		s.Require().NoError(err)
		s.NotNil(res.ResolveAt)
	})

	s.Run("resolve problem", func() {
		clientID := types.NewUserID()
		managerID := types.NewUserID()
		managerID2 := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		problem1, err := s.Database.Problem(s.Ctx).Create().
			SetChatID(chat.ID).
			SetManagerID(managerID2).
			Save(s.Ctx)
		s.Require().NoError(err)
		s.NotEmpty(problem1)

		res, err := s.Database.Problem(s.Ctx).Get(s.Ctx, problem1.ID)
		s.Require().NoError(err)
		s.Empty(res.ResolveAt)
		s.Equal(managerID2, res.ManagerID)

		err = s.repo.ResolveProblem(s.Ctx, problem1.ID, managerID)
		s.Require().Error(err)
	})
}
