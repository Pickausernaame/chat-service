package server

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
	eventstream "github.com/Pickausernaame/chat-service/internal/services/event-stream"
	"github.com/Pickausernaame/chat-service/internal/types"
	"github.com/Pickausernaame/chat-service/internal/validator"
	websocketstream "github.com/Pickausernaame/chat-service/internal/websocket-stream"
)

const (
	readHeaderTimeout = time.Second
	shutdownTimeout   = 3 * time.Second
)

type eventStream interface {
	Subscribe(ctx context.Context, userID types.UserID) (<-chan eventstream.Event, error)
}

type adapter interface {
	Adapt(event eventstream.Event) (any, error)
}

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options
type Options struct {
	serverName      string                 `option:"mandatory" validate:"required"`
	addr            string                 `option:"mandatory" validate:"required,hostname_port"`
	allowOrigins    []string               `option:"mandatory" validate:"min=1"`
	v1Swagger       *openapi3.T            `option:"mandatory" validate:"required"`
	reg             func(v *echo.Group)    `option:"mandatory" validate:"required"`
	keycloakClient  *keycloakclient.Client `option:"mandatory" validate:"required"`
	resource        string                 `option:"mandatory" validate:"required"`
	role            string                 `option:"mandatory" validate:"required"`
	secWsProtocol   string                 `option:"mandatory" validate:"required"`
	eventSubscriber eventStream            `option:"mandatory" validate:"required"`
	errHandler      echo.HTTPErrorHandler  `option:"mandatory" validate:"required"`
	adapter         adapter                `option:"mandatory" validate:"required"`
}

type Server struct {
	lg         *zap.Logger
	srv        *http.Server
	shutdownCh chan struct{}
}

func New(opts Options) (*Server, error) {
	if err := validator.Validator.Struct(opts); err != nil {
		return nil, fmt.Errorf("options validation error: %v", err)
	}

	e := echo.New()
	e.HTTPErrorHandler = opts.errHandler

	e.Use(
		middleware.RecoverWithConfig(middleware.RecoverConfig{
			Skipper:           middleware.DefaultSkipper,
			StackSize:         4 << 10,
			DisableStackAll:   false,
			DisablePrintStack: false,
			LogErrorFunc:      middlewares.RecoveryLogFunc,
		}),
		middlewares.ZapLogger(zap.L().Named(opts.serverName+" middleware")),
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
	shutdownCh := make(chan struct{})

	wsHandler, err := websocketstream.NewHTTPHandler(
		websocketstream.NewOptions(
			opts.eventSubscriber,
			opts.adapter,
			websocketstream.JSONEventWriter{},
			websocketstream.NewUpgrader(opts.allowOrigins, opts.secWsProtocol),
			shutdownCh))
	if err != nil {
		return nil, fmt.Errorf("making ws handler: %v", err)
	}

	e.GET("/ws", wsHandler.Serve)

	opts.reg(v1)

	return &Server{
		lg: zap.L().Named(opts.serverName),
		srv: &http.Server{
			Addr:              opts.addr,
			Handler:           e,
			ReadHeaderTimeout: readHeaderTimeout,
		},
		shutdownCh: shutdownCh,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	s.srv.RegisterOnShutdown(func() {
		close(s.shutdownCh)
	})

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
