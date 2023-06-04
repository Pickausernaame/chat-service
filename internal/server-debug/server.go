package serverdebug

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo-contrib/pprof"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/Pickausernaame/chat-service/internal/buildinfo"
	"github.com/Pickausernaame/chat-service/internal/logger"
	"github.com/Pickausernaame/chat-service/internal/middlewares"
	clientevents "github.com/Pickausernaame/chat-service/internal/server-client/events"
	clientv1 "github.com/Pickausernaame/chat-service/internal/server-client/v1"
	managerevents "github.com/Pickausernaame/chat-service/internal/server-manager/events"
	managerv1 "github.com/Pickausernaame/chat-service/internal/server-manager/v1"
)

const (
	readHeaderTimeout = time.Second
	shutdownTimeout   = 3 * time.Second
)

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options
type Options struct {
	addr string `option:"mandatory" validate:"required,hostname_port"`
}

type Server struct {
	lg  *zap.Logger
	srv *http.Server
}

func New(opts Options) (*Server, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validation debug server options error: %v", err)
	}

	lg := zap.L().Named("server-debug")

	e := echo.New()
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		Skipper:           middleware.DefaultSkipper,
		StackSize:         4 << 10,
		DisableStackAll:   false,
		DisablePrintStack: false,
		LogErrorFunc:      middlewares.RecoveryLogFunc,
	}),
		middlewares.ZapLogger(lg.Named("middleware")))

	s := &Server{
		lg: lg,
		srv: &http.Server{
			Addr:              opts.addr,
			Handler:           e,
			ReadHeaderTimeout: readHeaderTimeout,
		},
	}

	index := newIndexPage()
	index.addPage("/version", "Get build information")
	index.addPage("/debug/pprof/", "Go std profiler")
	index.addPage("/debug/pprof/profile?seconds=30", "Take half min profile")
	index.addPage("/debug/error", "Debug sentry error event")
	index.addPage("/schema/client", "Get client openAPI specification")
	index.addPage("/schema/manager", "Get client openAPI specification")
	index.addPage("/schema/client-events", "Get client events openAPI specification")
	index.addPage("/schema/manager-events", "Get manager events openAPI specification")

	e.GET("/", index.handler)
	e.GET("/version", s.Version)
	e.GET("/log/level", s.getLogLevel)
	e.PUT("/log/level", s.setLogLevel)
	e.GET("/debug/error", s.error)
	e.GET("/schema/client", s.clientSchema)
	e.GET("/schema/manager", s.managerSchema)
	e.GET("/schema/client-events", s.clientEventsSchema)
	e.GET("/schema/manager-events", s.managerEventsSchema)

	pprof.Register(e, "/debug/pprof")

	return s, nil
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

func (s *Server) Version(eCtx echo.Context) error {
	data, err := json.Marshal(buildinfo.BuildInfo)
	if err != nil {
		return eCtx.NoContent(http.StatusInternalServerError)
	}
	return eCtx.JSONBlob(http.StatusOK, data)
}

func (s *Server) getLogLevel(eCtx echo.Context) error {
	s.lg.Debug("getting log level")
	return eCtx.String(http.StatusOK, logger.LogLevel())
}

func (s *Server) setLogLevel(eCtx echo.Context) error {
	lvl := eCtx.FormValue("level")

	err := logger.SetLogLevel(logger.NewLogLevelOption(strings.ToLower(lvl)))
	s.lg.Debug("setting log level", zap.String("level", lvl), zap.Error(err))
	if err != nil {
		return eCtx.NoContent(http.StatusBadRequest)
	}

	return eCtx.String(http.StatusOK, s.lg.Level().String())
}

func (s *Server) error(eCtx echo.Context) error {
	s.lg.Debug("debug error msg", zap.Error(errors.New("debug error")))
	return eCtx.NoContent(http.StatusOK)
}

func (s *Server) clientSchema(eCtx echo.Context) error {
	spec, err := clientv1.GetSwagger()
	if err != nil {
		eCtx.Error(err)
	}
	return eCtx.JSON(http.StatusOK, spec)
}

func (s *Server) managerSchema(eCtx echo.Context) error {
	spec, err := managerv1.GetSwagger()
	if err != nil {
		eCtx.Error(err)
	}
	return eCtx.JSON(http.StatusOK, spec)
}

func (s *Server) clientEventsSchema(eCtx echo.Context) error {
	spec, err := clientevents.GetSwagger()
	if err != nil {
		eCtx.Error(err)
	}
	return eCtx.JSON(http.StatusOK, spec)
}

func (s *Server) managerEventsSchema(eCtx echo.Context) error {
	spec, err := managerevents.GetSwagger()
	if err != nil {
		eCtx.Error(err)
	}
	return eCtx.JSON(http.StatusOK, spec)
}
