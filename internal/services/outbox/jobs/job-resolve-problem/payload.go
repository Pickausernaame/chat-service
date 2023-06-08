package jobresolveproblem

import (
	"encoding/json"

	"github.com/Pickausernaame/chat-service/internal/types"
	"github.com/Pickausernaame/chat-service/internal/validator"
)

type Request struct {
	ChatID    types.ChatID    `json:"chat_id" validate:"required"`
	ManagerID types.UserID    `json:"manager_id" validate:"required"`
	MessageID types.MessageID `json:"message_id" validate:"required"`
	RequestID types.RequestID `json:"request_id" validate:"required"`
}

func MarshalPayload(managerID types.UserID, requestID types.RequestID, messageID types.MessageID, chatID types.ChatID) (string, error) {
	r := &Request{
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
