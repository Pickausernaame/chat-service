package sendclientmessagejob

import (
	"encoding/json"

	"github.com/Pickausernaame/chat-service/internal/types"
)

func MarshalPayload(messageID types.MessageID) (string, error) {
	if err := messageID.Validate(); err != nil {
		return "", err
	}
	res, err := json.Marshal(messageID)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
