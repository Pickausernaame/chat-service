package managerv1_test

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	internalerrors "github.com/Pickausernaame/chat-service/internal/errors"
	managerv1 "github.com/Pickausernaame/chat-service/internal/server-manager/v1"
	"github.com/Pickausernaame/chat-service/internal/types"
	getchathistory "github.com/Pickausernaame/chat-service/internal/usecases/manager/get-chat-history"
)

func (s *HandlersSuite) TestGetHistory_BindRequestError() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getHistory", `{"page_size":`)

	// Action.
	err := s.handlers.PostGetChatHistory(eCtx, managerv1.PostGetChatHistoryParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Equal(http.StatusBadRequest, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestGetHistory_Usecase_InvalidRequest() {
	// Arrange.
	reqID := types.NewRequestID()

	resp, eCtx := s.newEchoCtx(reqID, "/v1/getChatHistory",
		fmt.Sprintf(`{"chatId":"%s","pageSize":9}`, s.chatID.String()))
	s.getChatHistoryUseCase.EXPECT().Handle(eCtx.Request().Context(), getchathistory.Request{
		ID:        reqID,
		ManagerID: s.managerID,
		ChatID:    s.chatID,
		PageSize:  9,
		Cursor:    "",
	}).Return(getchathistory.Response{}, getchathistory.ErrInvalidRequest)

	// Action.
	err := s.handlers.PostGetChatHistory(eCtx, managerv1.PostGetChatHistoryParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Equal(http.StatusBadRequest, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestGetHistory_Usecase_InvalidCursor() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getChatHistory",
		fmt.Sprintf(`{"chatId":"%s","cursor":"abracadabra"}`, s.chatID.String()))

	s.getChatHistoryUseCase.EXPECT().Handle(eCtx.Request().Context(), getchathistory.Request{
		ID:        reqID,
		ManagerID: s.managerID,
		ChatID:    s.chatID,
		Cursor:    "abracadabra",
	}).Return(getchathistory.Response{}, getchathistory.ErrInvalidCursor)

	// Action.
	err := s.handlers.PostGetChatHistory(eCtx, managerv1.PostGetChatHistoryParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Equal(http.StatusBadRequest, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestGetHistory_Usecase_UnknownError() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getChatHistory",
		fmt.Sprintf(`{"chatId":"%s","pageSize":9}`, s.chatID.String()))
	s.getChatHistoryUseCase.EXPECT().Handle(eCtx.Request().Context(), getchathistory.Request{
		ID:        reqID,
		ManagerID: s.managerID,
		ChatID:    s.chatID,
		PageSize:  9,
		Cursor:    "",
	}).Return(getchathistory.Response{}, errors.New("something went wrong"))

	// Action.
	err := s.handlers.PostGetChatHistory(eCtx, managerv1.PostGetChatHistoryParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestGetHistory_Usecase_Success() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getChatHistory",
		fmt.Sprintf(`{"chatId":"%s","pageSize":10}`, s.chatID.String()))
	msgs := []*getchathistory.Message{
		{
			ID:        types.NewMessageID(),
			AuthorID:  types.NewUserID(),
			Body:      "hello!",
			CreatedAt: time.Unix(1, 1).UTC(),
		},
		{
			ID:        types.NewMessageID(),
			AuthorID:  types.NewUserID(),
			Body:      "service message",
			CreatedAt: time.Unix(2, 2).UTC(),
		},
	}

	s.getChatHistoryUseCase.EXPECT().Handle(eCtx.Request().Context(), getchathistory.Request{
		ID:        reqID,
		ManagerID: s.managerID,
		ChatID:    s.chatID,
		PageSize:  10,
	}).Return(getchathistory.Response{
		Messages:   msgs,
		NextCursor: "",
	}, nil)

	// Action.
	err := s.handlers.PostGetChatHistory(eCtx, managerv1.PostGetChatHistoryParams{XRequestID: reqID})

	// Assert.
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.Code)
	s.JSONEq(fmt.Sprintf(`
{
   "data":
   {
       "messages":
       [
           {
               "authorId": %q,
               "body": "hello!",
               "createdAt": "1970-01-01T00:00:01.000000001Z",
               "id": %q
           },
           {
               "authorId": %q,
               "body": "service message",
               "createdAt": "1970-01-01T00:00:02.000000002Z",
               "id": %q
           }
       ],
       "next": ""
   }
}`, msgs[0].AuthorID, msgs[0].ID, msgs[1].AuthorID, msgs[1].ID), resp.Body.String())
}
