package cursor

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func Encode(data any) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("marshaling: %v", err)
	}

	// Encode the JSON string to base64URL
	base64url := base64.URLEncoding.EncodeToString(jsonBytes)

	return base64url, nil
}

func Decode(in string, to any) error {
	jsonBytes, err := base64.URLEncoding.DecodeString(in)
	if err != nil {
		return fmt.Errorf("decoding base: %v", err)
	}

	if err = json.Unmarshal(jsonBytes, to); err != nil {
		return fmt.Errorf("unmarshaling: %v", err)
	}
	return nil
}
