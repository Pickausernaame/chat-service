package middlewares

import (
	"errors"

	"github.com/golang-jwt/jwt"

	keycloakclient "github.com/Pickausernaame/chat-service/internal/clients/keycloak"
	"github.com/Pickausernaame/chat-service/internal/types"
)

var (
	ErrNoAllowedResources = errors.New("no allowed resources")
	ErrSubjectNotDefined  = errors.New(`"sub" is not defined`)
)

//nolint:tagliatelle
type claims struct {
	Audience         keycloakclient.StringOrArray   `json:"aud,omitempty"`
	Subject          types.UserID                   `json:"sub,omitempty"`
	AllowedResources map[string]map[string][]string `json:"resource_access,omitempty"`
	jwt.StandardClaims
}

// Valid returns errors:
// - from StandardClaims validation;
// - ErrNoAllowedResources, if claims doesn't contain `resource_access` map or it's empty;
// - ErrSubjectNotDefined, if claims doesn't contain `sub` field or subject is zero UUID.
func (c claims) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}
	if c.Subject.IsZero() {
		return ErrSubjectNotDefined
	}
	if len(c.AllowedResources) == 0 {
		return ErrNoAllowedResources
	}

	return nil
}

func (c claims) UserID() types.UserID {
	return c.Subject
}
