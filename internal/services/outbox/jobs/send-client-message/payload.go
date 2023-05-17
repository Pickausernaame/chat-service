package sendclientmessagejob

import (
	"encoding/json"
	"errors"

	"github.com/Pickausernaame/chat-service/internal/types"
)

// FIXME: Вероятно необходимо добавить приватных типов и функций.

func MarshalPayload(messageID types.MessageID) (string, error) {
	if messageID == types.MessageIDNil {
		return "", errors.New("nil message id")
	}
	res, err := json.Marshal(messageID)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
