package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/Pickausernaame/chat-service/internal/config"
	"github.com/Pickausernaame/chat-service/internal/logger"
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

	isProduction := false
	v := os.Getenv("IS_PRODUCTION")
	if v != "" {
		isProduction, err = strconv.ParseBool(v)
		if err != nil {
			return fmt.Errorf("parse env error: %v", err)
		}
	}

	opts := logger.NewOptions(cfg.Log.Level, logger.WithProductionMode(isProduction))
	if err = logger.Init(opts); err != nil {
		return fmt.Errorf("logger init error: %v", err)
	}

	logger.Sync()

	srvDebug, err := serverdebug.New(serverdebug.NewOptions(cfg.Servers.Debug.Addr))
	if err != nil {
		return fmt.Errorf("init debug server: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	// Run servers.
	eg.Go(func() error { return srvDebug.Run(ctx) })

	// Run services.
	// Ждут своего часа.
	// ...

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("wait app stop: %v", err)
	}

	return nil
}
