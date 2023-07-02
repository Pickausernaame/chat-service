//go:build integration

package chatsrepo_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	chatsrepo "github.com/Pickausernaame/chat-service/internal/repositories/chats"
	"github.com/Pickausernaame/chat-service/internal/testingh"
	"github.com/Pickausernaame/chat-service/internal/types"
)

type ChatsRepoSuite struct {
	testingh.DBSuite
	repo *chatsrepo.Repo
}

func TestChatsRepoSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ChatsRepoSuite{DBSuite: testingh.NewDBSuite("TestChatsRepoSuite")})
}

func (s *ChatsRepoSuite) SetupSuite() {
	s.DBSuite.SetupSuite()

	var err error

	s.repo, err = chatsrepo.New(chatsrepo.NewOptions(s.Database))
	s.Require().NoError(err)
}

func (s *ChatsRepoSuite) Test_CreateIfNotExists() {
	s.Run("chat does not exist, should be created", func() {
		clientID := types.NewUserID()

		chatID, err := s.repo.CreateIfNotExists(s.Ctx, clientID)
		s.Require().NoError(err)
		s.NotEmpty(chatID)
	})

	s.Run("chat already exists", func() {
		clientID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		chatID, err := s.repo.CreateIfNotExists(s.Ctx, clientID)
		s.Require().NoError(err)
		s.Require().NotEmpty(chatID)
		s.Equal(chat.ID, chatID)
	})
}

func (s *ChatsRepoSuite) Test_ClientIDByID() {
	s.Run("get clientID by chatID - success", func() {
		clientID := types.NewUserID()

		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		res, err := s.repo.ClientIDByID(s.Ctx, chat.ID)
		s.Require().NoError(err)
		s.Equal(clientID, res)
	})

	s.Run("chat not exist", func() {
		clientID := types.NewUserID()
		chatID := types.NewChatID()

		_, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		res, err := s.repo.ClientIDByID(s.Ctx, chatID)
		s.Require().Error(err)
		s.Empty(res)
	})
}
