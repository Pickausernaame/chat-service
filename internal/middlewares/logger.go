package middlewares

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ZapLogger(log *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(eCtx echo.Context) error {
			start := time.Now()

			err := next(eCtx)
			if err != nil {
				eCtx.Error(err)
			}

			req := eCtx.Request()
			if req.Method == http.MethodOptions {
				return nil
			}
			res := eCtx.Response()

			fields := []zapcore.Field{
				zap.String("remote_ip", eCtx.RealIP()),
				zap.String("latency", time.Since(start).String()),
				zap.String("host", req.Host),
				zap.String("method", req.Method),
				zap.String("path", req.RequestURI),
				zap.String("user_agent", req.UserAgent()),
				zap.Int64("size", res.Size),
				zap.Int("status", res.Status),
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
