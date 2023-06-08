package managerv1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	srverr "github.com/Pickausernaame/chat-service/internal/errors"
	"github.com/Pickausernaame/chat-service/internal/middlewares"
	getchathistory "github.com/Pickausernaame/chat-service/internal/usecases/manager/get-chat-history"
	"github.com/Pickausernaame/chat-service/pkg/pointer"
)

func (h Handlers) PostGetChatHistory(eCtx echo.Context, params PostGetChatHistoryParams) error {
	ctx := eCtx.Request().Context()
	managerID := middlewares.MustUserID(eCtx)

	req := &GetChatHistoryRequest{}
	if err := eCtx.Bind(req); err != nil {
		return fmt.Errorf("binding GetHistory: %w", err)
	}

	resp, err := h.getChatHistoryUseCase.Handle(ctx, getchathistory.Request{
		ID:        params.XRequestID,
		ManagerID: managerID,
		ChatID:    req.ChatId,
		PageSize:  pointer.Indirect(req.PageSize),
		Cursor:    pointer.Indirect(req.Cursor),
	})
	if err != nil {
		if errors.Is(err, getchathistory.ErrInvalidCursor) {
			return srverr.NewServerError(http.StatusBadRequest, "invalid cursor", err)
		}

		if errors.Is(err, getchathistory.ErrInvalidRequest) {
			return srverr.NewServerError(http.StatusBadRequest, "invalid request", err)
		}
		return err
	}

	msgs := make([]Message, 0, len(resp.Messages))
	for _, m := range resp.Messages {
		msgs = append(msgs, toMessage(m))
	}

	return eCtx.JSON(http.StatusOK, &GetChatHistoryResponse{
		Data: &MessagesPage{Messages: msgs, Next: resp.NextCursor},
	})
}

func toMessage(msg *getchathistory.Message) Message {
	m := Message{
		Id:        msg.ID,
		Body:      msg.Body,
		CreatedAt: msg.CreatedAt,
		AuthorId:  msg.AuthorID,
	}
	return m
}
