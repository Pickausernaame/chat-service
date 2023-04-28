package clientv1

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/Pickausernaame/chat-service/internal/types"
)

var stub = MessagesPage{Messages: []Message{
	{
		AuthorId:  types.NewUserID(),
		Body:      "Здравствуйте! Разберёмся.",
		CreatedAt: time.Now(),
		Id:        types.NewMessageID(),
	},
	{
		AuthorId:  types.MustParse[types.UserID]("a391aae2-1a67-4288-be4b-f6eda6d50ffd"),
		Body:      "Привет! Не могу снять денег с карты,\nпишет 'карта заблокирована'",
		CreatedAt: time.Now().Add(-time.Minute),
		Id:        types.NewMessageID(),
	},
}}

func (h Handlers) PostGetHistory(eCtx echo.Context, params PostGetHistoryParams) error {
	h.logger.Info("PostGetHistory", zap.Any("params", params))
	return eCtx.JSON(http.StatusOK, map[string]interface{}{"data": stub})
}
