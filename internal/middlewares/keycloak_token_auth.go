package middlewares

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	keycloakclient "github.com/Pickausernaame/chat-service/internal/clients/keycloak"
	"github.com/Pickausernaame/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/introspector_mock.gen.go -package=middlewaresmocks Introspector

const (
	tokenCtxKey     = "user-token"
	websocketHeader = "Sec-WebSocket-Protocol"
)

var ErrNoRequiredResourceRole = errors.New("no required resource role")

type Introspector interface {
	IntrospectToken(ctx context.Context, token string) (*keycloakclient.IntrospectTokenResult, error)
}

// NewKeycloakTokenAuth returns a middleware that implements "active" authentication:
// each request is verified by the Keycloak server.
func NewKeycloakTokenAuth(introspector Introspector, resource, role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if v := c.Request().Header.Get(websocketHeader); v != "" {
				return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
					KeyLookup: "header:" + websocketHeader + ":chat-service-protocol",
					Validator: Validator(introspector, resource, role),
				})(next)(c)
			}
			return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
				KeyLookup:  "header:" + echo.HeaderAuthorization,
				AuthScheme: "Bearer",
				Validator:  Validator(introspector, resource, role),
			})(next)(c)
		}
	}
}

func Validator(introspector Introspector, resource, role string) middleware.KeyAuthValidator {
	return func(tokenStr string, eCtx echo.Context) (bool, error) {
		if v := eCtx.Request().Header.Get(websocketHeader); v != "" {
			tokenStr = tokenStr[2:]
		}
		res, err := introspector.IntrospectToken(eCtx.Request().Context(), tokenStr)
		if err != nil {
			return false, err
		}

		if !res.Active {
			return false, errors.New("token is inactive")
		}

		cl := &claims{}
		token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, cl)
		if err != nil {
			return false, fmt.Errorf("parsing token error: %v", err)
		}

		if err := cl.Valid(); err != nil {
			return false, err
		}

		data, ok := cl.AllowedResources[resource]
		if !ok {
			return false, ErrNoRequiredResourceRole
		}

		roles, ok := data["roles"]
		if !ok {
			return false, ErrNoRequiredResourceRole
		}
		ok = false
		for _, r := range roles {
			if r == role {
				ok = true
				break
			}
		}

		if !ok {
			return false, ErrNoRequiredResourceRole
		}

		eCtx.Set(tokenCtxKey, token)
		return true, nil
	}
}

func MustUserID(eCtx echo.Context) types.UserID {
	uid, ok := userID(eCtx)
	if !ok {
		panic("no user token in request context")
	}
	return uid
}

func UserID(eCtx echo.Context) (types.UserID, bool) {
	return userID(eCtx)
}

func userID(eCtx echo.Context) (types.UserID, bool) {
	t := eCtx.Get(tokenCtxKey)
	if t == nil {
		return types.UserIDNil, false
	}

	tt, ok := t.(*jwt.Token)
	if !ok {
		return types.UserIDNil, false
	}

	userIDProvider, ok := tt.Claims.(interface{ UserID() types.UserID })
	if !ok {
		return types.UserIDNil, false
	}
	return userIDProvider.UserID(), true
}
