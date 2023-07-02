package managerv1_test

import (
	"errors"
	"fmt"

	srverr "github.com/Pickausernaame/chat-service/internal/errors"
	managerv1 "github.com/Pickausernaame/chat-service/internal/server-manager/v1"
	"github.com/Pickausernaame/chat-service/internal/types"
	resolveproblem "github.com/Pickausernaame/chat-service/internal/usecases/manager/resolve-problem"
)

func (s *HandlersSuite) TestResolveProblem_Usecase_SomeError() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/closeChat", fmt.Sprintf(`{"chatId":"%s"}`, s.chatID))
	s.resolveProblemUseCase.EXPECT().Handle(eCtx.Request().Context(), resolveproblem.Request{
		RequestID: reqID,
		ChatID:    s.chatID,
		ManagerID: s.managerID,
	}).Return(resolveproblem.Response{}, errors.New("something went wrong"))

	// Action.
	err := s.handlers.PostCloseChat(eCtx, managerv1.PostCloseChatParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestResolveProblem_Usecase_ProblemNotFoundError() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/closeChat", fmt.Sprintf(`{"chatId":"%s"}`, s.chatID))
	s.resolveProblemUseCase.EXPECT().Handle(eCtx.Request().Context(), resolveproblem.Request{
		RequestID: reqID,
		ChatID:    s.chatID,
		ManagerID: s.managerID,
	}).Return(resolveproblem.Response{}, resolveproblem.ErrProblemNotFound)

	// Action.
	err := s.handlers.PostCloseChat(eCtx, managerv1.PostCloseChatParams{XRequestID: reqID})

	// Assert.
	errsrv := &srverr.ServerError{}
	s.Require().Error(err)
	s.Require().ErrorAs(err, &errsrv)
	s.Require().Equal(int(managerv1.ErrorCodeProblemNotExistError), errsrv.Code)
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestResolveProblem_Usecase_Success() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/closeChat", fmt.Sprintf(`{"chatId":"%s"}`, s.chatID))
	s.resolveProblemUseCase.EXPECT().Handle(eCtx.Request().Context(), resolveproblem.Request{
		RequestID: reqID,
		ChatID:    s.chatID,
		ManagerID: s.managerID,
	}).Return(resolveproblem.Response{}, nil)

	// Action.
	err := s.handlers.PostCloseChat(eCtx, managerv1.PostCloseChatParams{XRequestID: reqID})

	// Assert.
	s.Require().NoError(err)
	s.NotEmpty(resp.Body)
}
