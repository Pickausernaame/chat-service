package clientv1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	srverr "github.com/Pickausernaame/chat-service/internal/errors"
	"github.com/Pickausernaame/chat-service/internal/middlewares"
	"github.com/Pickausernaame/chat-service/internal/types"
	gethistory "github.com/Pickausernaame/chat-service/internal/usecases/client/get-history"
	"github.com/Pickausernaame/chat-service/pkg/pointer"
)

func (h Handlers) PostGetHistory(eCtx echo.Context, params PostGetHistoryParams) error {
	ctx := eCtx.Request().Context()
	clientID := middlewares.MustUserID(eCtx)

	req := &GetHistoryRequest{}
	if err := eCtx.Bind(req); err != nil {
		return fmt.Errorf("binding GetHistory: %w", err)
	}

	resp, err := h.getHistory.Handle(ctx, gethistory.Request{
		ID:       params.XRequestID,
		ClientID: clientID,
		PageSize: pointer.Indirect(req.PageSize),
		Cursor:   pointer.Indirect(req.Cursor),
	})
	if err != nil {
		if errors.Is(err, gethistory.ErrInvalidCursor) {
			return srverr.NewServerError(http.StatusBadRequest, "invalid cursor", err)
		}

		if errors.Is(err, gethistory.ErrInvalidRequest) {
			return srverr.NewServerError(http.StatusBadRequest, "invalid request", err)
		}
		return err
	}

	msgs := make([]Message, 0, len(resp.Messages))
	for _, m := range resp.Messages {
		msgs = append(msgs, toMessage(m))
	}

	return eCtx.JSON(http.StatusOK, &GetHistoryResponse{Data: &MessagesPage{Messages: msgs, Next: resp.NextCursor}})
}

func toMessage(msg *gethistory.Message) Message {
	m := Message{
		Id:         msg.ID,
		Body:       msg.Body,
		CreatedAt:  msg.CreatedAt,
		IsBlocked:  msg.IsBlocked,
		IsReceived: msg.IsReceived,
		IsService:  msg.IsService,
	}
	if !msg.IsService {
		m.AuthorId = pointer.PtrWithZeroAsNil[types.UserID](msg.AuthorID)
	}
	return m
}
