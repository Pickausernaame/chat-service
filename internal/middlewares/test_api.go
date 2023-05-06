package middlewares

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"github.com/Pickausernaame/chat-service/internal/types"
)

func SetToken(c echo.Context, uid types.UserID) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims{
		Subject:        uid,
		StandardClaims: jwt.StandardClaims{},
	})
	c.Set(tokenCtxKey, token)
}
