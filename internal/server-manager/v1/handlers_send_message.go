package managerv1

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Pickausernaame/chat-service/internal/middlewares"
	sendmessage "github.com/Pickausernaame/chat-service/internal/usecases/manager/send-message"
)

func (h Handlers) PostSendMessage(eCtx echo.Context, params PostSendMessageParams) error {
	ctx := eCtx.Request().Context()
	managerID := middlewares.MustUserID(eCtx)

	req := &SendMessageRequest{}
	if err := eCtx.Bind(req); err != nil {
		return fmt.Errorf("binding SendMessageRequest: %w", err)
	}

	resp, err := h.sendMessageUseCase.Handle(ctx, sendmessage.Request{
		ID:          params.XRequestID,
		ManagerID:   managerID,
		ChatID:      req.ChatId,
		MessageBody: req.MessageBody,
	})
	if err != nil {
		return fmt.Errorf("handling send message usecase: %w", err)
	}

	return eCtx.JSON(http.StatusOK, SendMessageResponse{Data: &MessageWithoutBody{
		AuthorId:  managerID,
		CreatedAt: resp.CreatedAt,
		Id:        resp.MessageID,
	}})
}
