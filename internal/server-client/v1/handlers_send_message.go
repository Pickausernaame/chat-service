package clientv1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	errsrv "github.com/Pickausernaame/chat-service/internal/errors"
	"github.com/Pickausernaame/chat-service/internal/middlewares"
	sendmessage "github.com/Pickausernaame/chat-service/internal/usecases/client/send-message"
	"github.com/Pickausernaame/chat-service/pkg/pointer"
)

func (h Handlers) PostSendMessage(eCtx echo.Context, params PostSendMessageParams) error {
	ctx := eCtx.Request().Context()
	clientID := middlewares.MustUserID(eCtx)

	req := &SendMessageRequest{}
	if err := eCtx.Bind(req); err != nil {
		return fmt.Errorf("binding SendMessageRequest: %w", err)
	}

	r := sendmessage.Request{
		ID:          params.XRequestID,
		ClientID:    clientID,
		MessageBody: req.MessageBody,
	}

	resp, err := h.sendMessage.Handle(ctx, r)
	if err != nil {
		if errors.Is(err, sendmessage.ErrChatNotCreated) {
			return errsrv.NewServerError(int(ErrorCodeCreateChatError), "chat not created", err)
		}

		if errors.Is(err, sendmessage.ErrProblemNotCreated) {
			return errsrv.NewServerError(int(ErrorCodeCreateProblemError), "problem not created", err)
		}

		if errors.Is(err, sendmessage.ErrInvalidRequest) {
			return errsrv.NewServerError(http.StatusBadRequest, "invalid request", err)
		}

		return fmt.Errorf("handle sendMessage: %w", err)
	}

	return eCtx.JSON(http.StatusOK, SendMessageResponse{
		Data: &MessageHeader{
			AuthorId:  pointer.Ptr(resp.AuthorID),
			CreatedAt: pointer.Ptr(resp.CreatedAt),
			Id:        pointer.Ptr(resp.MessageID),
		},
	})
}
