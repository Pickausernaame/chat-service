package errhandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/Pickausernaame/chat-service/internal/errors"
)

var _ echo.HTTPErrorHandler = Handler{}.Handle

//go:generate options-gen -out-filename=errhandler_options.gen.go -from-struct=Options
type Options struct {
	logger          *zap.Logger                                    `option:"mandatory" validate:"required"`
	productionMode  bool                                           `option:"mandatory"`
	responseBuilder func(code int, msg string, details string) any `option:"mandatory" validate:"required"`
}

type Handler struct {
	lg              *zap.Logger
	productionMode  bool
	responseBuilder func(code int, msg string, details string) any
}

func New(opts Options) (Handler, error) {
	if err := opts.Validate(); err != nil {
		return Handler{}, err
	}
	return Handler{lg: opts.logger, productionMode: opts.productionMode, responseBuilder: opts.responseBuilder}, nil
}

func (h Handler) Handle(err error, eCtx echo.Context) {
	code, msg, details := errors.ProcessServerError(err)
	if h.productionMode {
		details = ""
	}

	err = eCtx.JSON(http.StatusOK, h.responseBuilder(code, msg, details))
	if err != nil {
		h.lg.Error("creating response error", zap.Error(err), zap.Int("code", code),
			zap.String("msg", msg), zap.String("details", details))
	}
}
