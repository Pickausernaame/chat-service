package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	keycloakclient "github.com/Pickausernaame/chat-service/internal/clients/keycloak"
	"github.com/Pickausernaame/chat-service/internal/config"
	"github.com/Pickausernaame/chat-service/internal/logger"
	chatsrepo "github.com/Pickausernaame/chat-service/internal/repositories/chats"
	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	problemsrepo "github.com/Pickausernaame/chat-service/internal/repositories/problems"
	serverdebug "github.com/Pickausernaame/chat-service/internal/server-debug"
	"github.com/Pickausernaame/chat-service/internal/store"
)

var configPath = flag.String("с", "configs/config.toml", "Path to config file")

func main() {
	if err := run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}

func run() (errReturned error) {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// parsing config
	cfg, err := config.ParseAndValidate(*configPath)
	if err != nil {
		return fmt.Errorf("parse and validate config %q: %v", *configPath, err)
	}

	// initialization logger
	if err = initLogger(cfg); err != nil {
		return fmt.Errorf("logger init error: %v", err)
	}
	defer logger.Sync()

	// initialization keycloak client
	kc, err := initKeycloak(cfg)
	if err != nil {
		return fmt.Errorf("keycloak client init error: %v", err)
	}

	// initialization psql client
	pg, err := initPSQLClient(cfg.Clients.PSQL)
	if err != nil {
		return fmt.Errorf("psql client init error: %v", err)
	}
	defer pg.Close()

	// run migrations
	if err = pg.Schema.Create(ctx); err != nil {
		return fmt.Errorf("psql migration error: %v", err)
	}

	// init database
	db := store.NewDatabase(pg)

	// creating repos
	// initialization messages repo
	msgRepo, err := messagesrepo.New(messagesrepo.NewOptions(db))
	if err != nil {
		return fmt.Errorf("init messages repo error: %v", err)
	}

	// initialization chat repo
	chatRepo, err := chatsrepo.New(chatsrepo.NewOptions(db))
	if err != nil {
		return fmt.Errorf("init chat repo error: %v", err)
	}

	// initialization problem repo
	problemRepo, err := problemsrepo.New(problemsrepo.NewOptions(db))
	if err != nil {
		return fmt.Errorf("init problem repo error: %v", err)
	}

	// ... other repos

	// creating servers
	// initialization debug server
	srvDebug, err := serverdebug.New(serverdebug.NewOptions(cfg.Servers.Debug.Addr))
	if err != nil {
		return fmt.Errorf("init debug server: %v", err)
	}

	// initialization client server
	srvClient, err := initServerClient(cfg, kc, msgRepo, chatRepo, problemRepo, db)
	if err != nil {
		return fmt.Errorf("init server client: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)
	// Run servers.
	// debug server
	eg.Go(func() error { return srvDebug.Run(ctx) })

	// server client
	eg.Go(func() error { return srvClient.Run(ctx) })

	// Run services.
	// Ждут своего часа.
	// ...

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("wait app stop: %v", err)
	}

	return nil
}

func initLogger(cfg config.Config) error {
	opts := logger.NewOptions(
		cfg.Global.Env,
		logger.WithLevel(logger.NewLogLevelOption(cfg.Log.Level)),
		logger.WithVersion(cfg.Global.Version),
		logger.WithSentryDSN(cfg.Sentry.DSN))

	return logger.Init(opts)
}

func initKeycloak(cfg config.Config) (*keycloakclient.Client, error) {
	if cfg.Clients.Keycloak.DebugMode && cfg.Global.Env == "prod" {
		zap.L().Warn("Keycloak.DebugMode = true AND Global.Env = Production")
	}

	return keycloakclient.New(
		keycloakclient.NewOptions(
			cfg.Clients.Keycloak.BasePath,
			cfg.Clients.Keycloak.Realm,
			cfg.Clients.Keycloak.ClientID,
			cfg.Clients.Keycloak.ClientSecret,
			keycloakclient.WithDebugMode(cfg.Clients.Keycloak.DebugMode)))
}

func initPSQLClient(cfg config.PSQLClientConfig) (*store.Client, error) {
	return store.NewPSQLClient(
		store.NewPSQLOptions(
			cfg.Host,
			cfg.UserName,
			cfg.Password,
			cfg.DBName,
			store.WithDebug(cfg.Debug)))
}
