package managerassignedtoproblemjob

import (
	"encoding/json"

	"github.com/Pickausernaame/chat-service/internal/types"
	"github.com/Pickausernaame/chat-service/internal/validator"
)

type request struct {
	ClientID  types.UserID    `validate:"required"`
	ManagerID types.UserID    `validate:"required"`
	RequestID types.RequestID `validate:"required"`
	MessageID types.MessageID `validate:"required"`
}

func MarshalPayload(clientID types.UserID, managerID types.UserID, requestID types.RequestID, messageID types.MessageID) (string, error) {
	r := &request{
		ClientID:  clientID,
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
