package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Pickausernaame/chat-service/internal/errors"
)

func ZapLogger(log *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(eCtx echo.Context) error {
			start := time.Now()
			code := -1
			err := next(eCtx)
			if err != nil {
				code = errors.GetServerErrorCode(err)
				eCtx.Error(err)
			}

			req := eCtx.Request()
			if req.Method == http.MethodOptions {
				return nil
			}
			res := eCtx.Response()
			if code == -1 {
				code = res.Status
			}

			fields := []zapcore.Field{
				zap.String("remote_ip", eCtx.RealIP()),
				zap.String("latency", time.Since(start).String()),
				zap.String("host", req.Host),
				zap.String("method", req.Method),
				zap.String("path", req.RequestURI),
				zap.String("user_agent", req.UserAgent()),
				zap.Int64("size", res.Size),
				zap.Int("status", code),
			}

			uID, ok := UserID(eCtx)
			if ok {
				fields = append(fields, zap.String("user_id", uID.String()))
			} else {
				fields = append(fields, zap.String("user_id", ""))
			}

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}
			fields = append(fields, zap.String("request_id", id))

			switch n := res.Status; {
			case n >= 500:
				log.With(zap.Error(err)).Error("Server error", fields...)
			case n >= 400:
				log.With(zap.Error(err)).Warn("Client error", fields...)
			case n >= 300:
				log.Info("Redirection", fields...)
			default:
				log.Info("Success", fields...)
			}

			return nil
		}
	}
}

func RecoveryLogFunc(c echo.Context, err error, stack []byte) error {
	l := zap.L().Named("recovery")
	msg := fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack)
	//nolint:exhaustive
	switch l.Level() {
	case zap.DebugLevel:
		c.Logger().Debug(msg)
	case zap.InfoLevel:
		c.Logger().Info(msg)
	case zap.WarnLevel:
		c.Logger().Warn(msg)
	case zap.ErrorLevel:
		c.Logger().Error(msg)
	default:
		c.Logger().Print(msg)
	}
	return nil
}
