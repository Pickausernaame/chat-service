package jobresolveproblem

import (
	"encoding/json"

	"github.com/Pickausernaame/chat-service/internal/types"
	"github.com/Pickausernaame/chat-service/internal/validator"
)

type request struct {
	ChatID    types.ChatID    `validate:"required"`
	ManagerID types.UserID    `validate:"required"`
	MessageID types.MessageID `validate:"required"`
	RequestID types.RequestID `validate:"required"`
}

func MarshalPayload(managerID types.UserID, requestID types.RequestID,
	messageID types.MessageID, chatID types.ChatID,
) (string, error) {
	r := &request{
		ChatID:    chatID,
		ManagerID: managerID,
		RequestID: requestID,
		MessageID: messageID,
	}
	if err := validator.Validator.Struct(r); err != nil {
		return "", err
	}
	res, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
