package managerv1_test

import (
	"errors"
	"fmt"
	"strings"

	managerv1 "github.com/Pickausernaame/chat-service/internal/server-manager/v1"
	"github.com/Pickausernaame/chat-service/internal/types"
	getassignedproblems "github.com/Pickausernaame/chat-service/internal/usecases/manager/get-assigned-problems"
)

func (s *HandlersSuite) TestGetChats_Usecase_Error() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getChats", "")
	s.getAssignedProblemsUseCase.EXPECT().Handle(eCtx.Request().Context(), getassignedproblems.Request{
		ManagerID: s.managerID,
	}).Return(getassignedproblems.Response{}, errors.New("something went wrong"))

	// Action.
	err := s.handlers.PostGetChats(eCtx, managerv1.PostGetChatsParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestGetChats_Success() {
	// Arrange.

	reqID := types.NewRequestID()
	chats := make([]*getassignedproblems.Chat, 0, 10)
	jsonChats := make([]string, 0, 10)
	template := `{
"chatId": "%s",
"clientId": "%s"
}`
	for i := 0; i < 10; i++ {
		chatID := types.NewChatID()
		clientID := types.NewUserID()
		chats = append(chats, &getassignedproblems.Chat{
			ChatID:   chatID,
			ClientID: clientID,
		})
		jsonChats = append(jsonChats, fmt.Sprintf(template, chatID.String(), clientID.String()))
	}

	resp, eCtx := s.newEchoCtx(reqID, "/v1/getChats", "")
	s.getAssignedProblemsUseCase.EXPECT().Handle(eCtx.Request().Context(), getassignedproblems.Request{
		ManagerID: s.managerID,
	}).Return(getassignedproblems.Response{Chats: chats}, nil)

	// Action.
	err := s.handlers.PostGetChats(eCtx, managerv1.PostGetChatsParams{XRequestID: reqID})

	// Assert.
	s.Require().NoError(err)

	expectedJson := fmt.Sprintf("{\n    \"data\":\n    {\n        \"chats\":[%s]}}", strings.Join(jsonChats, ","))
	s.JSONEq(expectedJson, resp.Body.String())
}
