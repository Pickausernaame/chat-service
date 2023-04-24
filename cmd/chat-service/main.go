package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/Pickausernaame/chat-service/internal/config"
	"github.com/Pickausernaame/chat-service/internal/logger"
	clientv1 "github.com/Pickausernaame/chat-service/internal/server-client/v1"
	serverdebug "github.com/Pickausernaame/chat-service/internal/server-debug"
)

var configPath = flag.String("config", "configs/config.toml", "Path to config file")

func main() {
	if err := run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}

func run() (errReturned error) {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.ParseAndValidate(*configPath)
	if err != nil {
		return fmt.Errorf("parse and validate config %q: %v", *configPath, err)
	}

	opts := logger.NewOptions(cfg.Global.Env, logger.WithLevel(logger.NewLogLevelOption(cfg.Log.Level)),
		logger.WithVersion(cfg.Global.Version), logger.WithSentryDSN(cfg.Sentry.DSN))
	if err = logger.Init(opts); err != nil {
		return fmt.Errorf("logger init error: %v", err)
	}

	defer logger.Sync()

	srvDebug, err := serverdebug.New(serverdebug.NewOptions(cfg.Servers.Debug.Addr))
	if err != nil {
		return fmt.Errorf("init debug server: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	// Run servers.
	// debug server
	eg.Go(func() error { return srvDebug.Run(ctx) })

	swg, err := clientv1.GetSwagger()
	if err != nil {
		return fmt.Errorf("getting swagger: %v", err)
	}
	srvClient, err := initServerClient(cfg.Servers.Client.Addr, cfg.Servers.Client.AllowsOrigins, swg)
	if err != nil {
		return fmt.Errorf("init server client: %v", err)
	}
	// client
	eg.Go(func() error { return srvClient.Run(ctx) })

	// Run services.
	// Ждут своего часа.
	// ...

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("wait app stop: %v", err)
	}

	return nil
}
