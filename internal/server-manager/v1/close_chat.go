package managerv1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	srverr "github.com/Pickausernaame/chat-service/internal/errors"
	"github.com/Pickausernaame/chat-service/internal/middlewares"
	resolveproblem "github.com/Pickausernaame/chat-service/internal/usecases/manager/resolve-problem"
)

func (h Handlers) PostCloseChat(eCtx echo.Context, params PostCloseChatParams) error {
	ctx := eCtx.Request().Context()
	managerID := middlewares.MustUserID(eCtx)

	req := &ChatId{}
	if err := eCtx.Bind(req); err != nil {
		return fmt.Errorf("binding GetHistory: %w", err)
	}

	_, err := h.resolveProblemUseCase.Handle(ctx, resolveproblem.Request{
		ChatID:    req.ChatId,
		ManagerID: managerID,
		RequestID: params.XRequestID,
	})
	if err != nil {
		if errors.Is(err, resolveproblem.ErrProblemNotFound) {
			return srverr.NewServerError(int(ErrorCodeProblemNotExistError), "problem not found", err)
		}

		if errors.Is(err, resolveproblem.ErrInvalidRequest) {
			return srverr.NewServerError(http.StatusBadRequest, "invalid request", err)
		}
		return fmt.Errorf("resolve problem: %v", err)
	}

	return eCtx.JSON(http.StatusOK, CloseChatResponse{})
}
