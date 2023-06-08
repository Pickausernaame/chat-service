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
	jobsrepo "github.com/Pickausernaame/chat-service/internal/repositories/jobs"
	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	problemsrepo "github.com/Pickausernaame/chat-service/internal/repositories/problems"
	serverdebug "github.com/Pickausernaame/chat-service/internal/server-debug"
	afcverdictsprocessor "github.com/Pickausernaame/chat-service/internal/services/afc-verdicts-processor"
	eventstream "github.com/Pickausernaame/chat-service/internal/services/event-stream"
	inmemeventstream "github.com/Pickausernaame/chat-service/internal/services/event-stream/in-mem"
	managerload "github.com/Pickausernaame/chat-service/internal/services/manager-load"
	inmemmanagerpool "github.com/Pickausernaame/chat-service/internal/services/manager-pool/in-mem"
	managerscheduler "github.com/Pickausernaame/chat-service/internal/services/manager-scheduler"
	msgproducer "github.com/Pickausernaame/chat-service/internal/services/msg-producer"
	"github.com/Pickausernaame/chat-service/internal/services/outbox"
	clientmessageblockedjob "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/client-message-blocked"
	clientmessagesentjob "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/client-message-sent"
	jobresolveproblem "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/job-resolve-problem"
	managerassignedtoproblemjob "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/manager-assigned-to-problem"
	sendclientmessagejob "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/send-client-message"
	sendmanagermessagejob "github.com/Pickausernaame/chat-service/internal/services/outbox/jobs/send-manager-message"
	"github.com/Pickausernaame/chat-service/internal/store"
	"github.com/Pickausernaame/chat-service/internal/types"
)

var configPath = flag.String("с", "configs/config.toml", "Path to config file")

func main() {
	if err := run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}

type eventSubscriber interface {
	Subscribe(ctx context.Context, userID types.UserID) (<-chan eventstream.Event, error)
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

	// initialization job repo
	jobRepo, err := jobsrepo.New(jobsrepo.NewOptions(db))
	if err != nil {
		return fmt.Errorf("init job repo error: %v", err)
	}

	// ... other repos

	// creating services
	// initialization msgProducer service
	// initialization kafka writer
	kw := msgproducer.NewKafkaWriter(cfg.Service.MsgSender.Brokers,
		cfg.Service.MsgSender.Topic, cfg.Service.MsgSender.BatchSize)

	managerKw := msgproducer.NewKafkaWriter(cfg.Service.MsgSender.Brokers,
		cfg.Service.MsgSender.Topic, cfg.Service.MsgSender.BatchSize)

	msgProdService, err := msgproducer.New(
		msgproducer.NewOptions(kw, msgproducer.WithEncryptKey(cfg.Service.MsgSender.EncryptionKey)))
	if err != nil {
		return fmt.Errorf("init msg sender service: %v", err)
	}

	managerMsgProdService, err := msgproducer.New(
		msgproducer.NewOptions(managerKw, msgproducer.WithEncryptKey(cfg.Service.MsgSender.EncryptionKey)))
	if err != nil {
		return fmt.Errorf("init msg sender service: %v", err)
	}

	// initialization outbox service
	obox, err := outbox.New(outbox.NewOptions(cfg.Service.Outbox.Workers, cfg.Service.Outbox.Idle,
		cfg.Service.Outbox.ReservedFor, jobRepo, db))
	if err != nil {
		return fmt.Errorf("init outbox service: %v", err)
	}

	eventStream := inmemeventstream.New()
	// initialization sendMsg job
	sendMsgJob, err := sendclientmessagejob.New(sendclientmessagejob.NewOptions(msgProdService, msgRepo, eventStream))
	if err != nil {
		return fmt.Errorf("init send msg job: %v", err)
	}

	err = obox.RegisterJob(sendMsgJob)
	if err != nil {
		return fmt.Errorf("registration send msg job: %v", err)
	}

	sendManagerMsgJob, err := sendmanagermessagejob.New(sendmanagermessagejob.NewOptions(managerMsgProdService, msgRepo, chatRepo, eventStream))
	if err != nil {
		return fmt.Errorf("init send manager msg job: %v", err)
	}

	err = obox.RegisterJob(sendManagerMsgJob)
	if err != nil {
		return fmt.Errorf("registration send manager msg job: %v", err)
	}

	// initialization sendMsg job
	msgBlockedJob, err := clientmessageblockedjob.New(clientmessageblockedjob.NewOptions(msgRepo, eventStream))
	if err != nil {
		return fmt.Errorf("init msg blocked job: %v", err)
	}

	err = obox.RegisterJob(msgBlockedJob)
	if err != nil {
		return fmt.Errorf("registration msg blocked job: %v", err)
	}

	msgSentJob, err := clientmessagesentjob.New(clientmessagesentjob.NewOptions(msgRepo, problemRepo, eventStream))
	if err != nil {
		return fmt.Errorf("init msg sent job: %v", err)
	}

	err = obox.RegisterJob(msgSentJob)
	if err != nil {
		return fmt.Errorf("registration msg sent job: %v", err)
	}

	// initialization manager pool service
	manPoolService := inmemmanagerpool.New()

	// initialization manager load service
	manLoadService, err := managerload.New(managerload.NewOptions(cfg.Service.ManagerLoad.MaxProblemsAtSameTime, problemRepo))
	if err != nil {
		return fmt.Errorf("init manLoadService service: %v", err)
	}

	mngrAssignedJob, err := managerassignedtoproblemjob.New(managerassignedtoproblemjob.NewOptions(msgRepo, manLoadService, eventStream))
	if err != nil {
		return fmt.Errorf("init manager assigned job: %v", err)
	}

	err = obox.RegisterJob(mngrAssignedJob)
	if err != nil {
		return fmt.Errorf("registration manager assigned job: %v", err)
	}

	resolveProblem, err := jobresolveproblem.New(jobresolveproblem.NewOptions(msgRepo, chatRepo, manLoadService, eventStream))
	if err != nil {
		return fmt.Errorf("init resolve problem job: %v", err)
	}

	err = obox.RegisterJob(resolveProblem)
	if err != nil {
		return fmt.Errorf("registration resolve problem job: %v", err)
	}

	mngrScheduler, err := managerscheduler.New(managerscheduler.NewOptions(cfg.Service.ManagerScheduler.Period, manPoolService, msgRepo, obox, problemRepo, db))
	if err != nil {
		return fmt.Errorf("init manager scheduler error: %v", err)
	}

	// initialization afc-verdicts-processor service
	afc := cfg.Service.AvcVerdictProcessor
	afcProcessor, err := afcverdictsprocessor.New(
		afcverdictsprocessor.NewOptions(afc.Brokers,
			afc.Consumers,
			afc.ConsumerGroup,
			afc.VerdictsTopic,
			afcverdictsprocessor.NewKafkaReader,
			afcverdictsprocessor.NewKafkaDLQWriter(afc.Brokers, afc.DlqTopic),
			db,
			msgRepo,
			obox,
			afcverdictsprocessor.WithVerdictsSignKey(afc.EncryptKey),
		))
	if err != nil {
		return fmt.Errorf("init afc-verdicts-processor service: %v", err)
	}

	// creating servers
	// initialization debug server
	srvDebug, err := serverdebug.New(serverdebug.NewOptions(cfg.Servers.Debug.Addr))
	if err != nil {
		return fmt.Errorf("init debug server: %v", err)
	}

	// initialization client server
	srvClient, err := initServerClient(cfg, kc, msgRepo, chatRepo, problemRepo, db, obox, eventStream)
	if err != nil {
		return fmt.Errorf("init server client: %v", err)
	}

	// initialization manager server
	srvManager, err := initServerManager(cfg, kc, manLoadService, manPoolService, eventStream, chatRepo,
		problemRepo, msgRepo, obox, db)
	if err != nil {
		return fmt.Errorf("init server manager: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)
	// Run servers.
	// debug server
	eg.Go(func() error { return srvDebug.Run(ctx) })

	// server client
	eg.Go(func() error { return srvClient.Run(ctx) })

	// server manager
	eg.Go(func() error { return srvManager.Run(ctx) })

	// Run services.
	// outbox run
	eg.Go(func() error { return obox.Run(ctx) })

	eg.Go(func() error { return afcProcessor.Run(ctx) })

	eg.Go(func() error { return mngrScheduler.Run(ctx) })

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("wait app stop: %v", err)
	}

	// close
	if err := msgProdService.Close(); err != nil {
		zap.L().Error("closing msgProdService error", zap.Error(err))
	}

	if err := kw.Close(); err != nil {
		zap.L().Error("closing kafka writer error", zap.Error(err))
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
