package managerv1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Pickausernaame/chat-service/internal/middlewares"
	getassignedproblems "github.com/Pickausernaame/chat-service/internal/usecases/manager/get-assigned-problems"
)

func (h Handlers) PostGetChats(eCtx echo.Context, params PostGetChatsParams) error {
	ctx := eCtx.Request().Context()
	managerID := middlewares.MustUserID(eCtx)

	req := getassignedproblems.Request{ManagerID: managerID}
	data, err := h.getAssignedProblemsUseCase.Handle(ctx, req)
	if err != nil {
		return err
	}

	res := make([]Chat, 0, len(data.Chats))
	for _, c := range data.Chats {
		res = append(res, Chat{
			ChatId:   c.ChatID,
			ClientId: c.ClientID,
		})
	}

	return eCtx.JSON(http.StatusOK, PostGetChats{Data: &ChatList{res}})
}
