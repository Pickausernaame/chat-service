package serverclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	oapimdlwr "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	keycloakclient "github.com/Pickausernaame/chat-service/internal/clients/keycloak"
	"github.com/Pickausernaame/chat-service/internal/middlewares"
	clientv1 "github.com/Pickausernaame/chat-service/internal/server-client/v1"
	"github.com/Pickausernaame/chat-service/internal/validator"
)

const (
	readHeaderTimeout = time.Second
	shutdownTimeout   = 3 * time.Second
)

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options
type Options struct {
	logger         *zap.Logger              `option:"mandatory" validate:"required"`
	addr           string                   `option:"mandatory" validate:"required,hostname_port"`
	allowOrigins   []string                 `option:"mandatory" validate:"min=1"`
	v1Swagger      *openapi3.T              `option:"mandatory" validate:"required"`
	v1Handlers     clientv1.ServerInterface `option:"mandatory" validate:"required"`
	keycloakClient *keycloakclient.Client   `option:"mandatory" validate:"required"`
	resource       string                   `option:"mandatory" validate:"required"`
	role           string                   `option:"mandatory" validate:"required"`
}

type Server struct {
	lg  *zap.Logger
	srv *http.Server
}

func New(opts Options) (*Server, error) {
	if err := validator.Validator.Struct(opts); err != nil {
		return nil, fmt.Errorf("options validation error: %v", err)
	}

	e := echo.New()
	e.Use(
		middleware.RecoverWithConfig(middleware.RecoverConfig{
			Skipper:           middleware.DefaultSkipper,
			StackSize:         4 << 10,
			DisableStackAll:   false,
			DisablePrintStack: false,
			LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
				l := zap.L().Named("recovery")
				msg := fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack)
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
			},
		}),
		middlewares.ZapLogger(opts.logger.Named("middleware")),
		// (165(size of message without body) + 3000*4(max size of body)) * 1(count of messages per 1 request) * 2 (
		// margin factor) --> 24Kb
		middleware.BodyLimit("24K"),
		middleware.CORSWithConfig(
			middleware.CORSConfig{
				AllowOrigins: opts.allowOrigins,
				AllowMethods: []string{"POST"},
			}),
		middlewares.NewKeycloakTokenAuth(opts.keycloakClient, opts.resource, opts.role),
	)

	v1 := e.Group("v1", oapimdlwr.OapiRequestValidatorWithOptions(opts.v1Swagger, &oapimdlwr.Options{
		Options: openapi3filter.Options{
			ExcludeRequestBody:  false,
			ExcludeResponseBody: true,
			AuthenticationFunc:  openapi3filter.NoopAuthenticationFunc,
		},
	}))

	clientv1.RegisterHandlers(v1, opts.v1Handlers)

	return &Server{
		lg: opts.logger,
		srv: &http.Server{
			Addr:              opts.addr,
			Handler:           e,
			ReadHeaderTimeout: readHeaderTimeout,
		},
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		return s.srv.Shutdown(ctx)
	})

	eg.Go(func() error {
		s.lg.Info("listen and serve", zap.String("addr", s.srv.Addr))
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("listen and serve: %v", err)
		}
		return nil
	})

	return eg.Wait()
}
