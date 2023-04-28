package keycloakclient

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type StringOrArray []string

type IntrospectTokenResult struct {
	Exp    int           `json:"exp"`
	Iat    int           `json:"iat"`
	Aud    StringOrArray `json:"aud"`
	Active bool          `json:"active"`
}

func (s *StringOrArray) UnmarshalJSON(data []byte) error {
	// array
	if len(data) > 1 && data[0] == '[' {
		var obj []string
		if err := json.Unmarshal(data, &obj); err != nil {
			return err
		}
		*s = obj
		return nil
	}

	// string
	var obj string
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	*s = []string{obj}
	return nil
}

// IntrospectToken implements
// https://www.keycloak.org/docs/latest/authorization_services/index.html#obtaining-information-about-an-rpt
func (c *Client) IntrospectToken(ctx context.Context, token string) (*IntrospectTokenResult, error) {
	url := fmt.Sprintf("realms/%s/protocol/openid-connect/token/introspect", c.realm)

	resp, err := c.R(ctx).SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"token_type_hint": "requesting_party_token",
			"token":           token,
		}).Post(url)
	if err != nil {
		return nil, fmt.Errorf("token introspect request error: %v", err)
	}
	res := &IntrospectTokenResult{}
	err = json.Unmarshal(resp.Body(), res)
	if err != nil {
		return nil, fmt.Errorf("token introspect unmarshal error: %v", err)
	}
	return res, nil
}

func (c *Client) R(ctx context.Context) *resty.Request {
	return c.cli.R().SetContext(ctx)
}
