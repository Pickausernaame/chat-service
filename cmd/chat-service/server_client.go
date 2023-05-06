package main

import (
	"fmt"

	"go.uber.org/zap"

	keycloakclient "github.com/Pickausernaame/chat-service/internal/clients/keycloak"
	"github.com/Pickausernaame/chat-service/internal/config"
	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	serverclient "github.com/Pickausernaame/chat-service/internal/server-client"
	"github.com/Pickausernaame/chat-service/internal/server-client/errhandler"
	clientv1 "github.com/Pickausernaame/chat-service/internal/server-client/v1"
	gethistory "github.com/Pickausernaame/chat-service/internal/usecases/client/get-history"
)

const nameServerClient = "server-client"

func initServerClient(
	cfg config.Config,
	keycloakClient *keycloakclient.Client,
	msgRepo *messagesrepo.Repo,
) (*serverclient.Server, error) {
	lg := zap.L().Named(nameServerClient)

	// getting specification
	v1Swagger, err := clientv1.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("getting swagger: %v", err)
	}

	// initialization errorHandler
	errHandler, err := errhandler.New(errhandler.NewOptions(lg, cfg.Global.IsProd(), errhandler.ResponseBuilder))
	if err != nil {
		return nil, fmt.Errorf("init errror handler: %v", err)
	}

	// creating useCases
	// initialization getHistory useCase
	getHistoryUC, err := gethistory.New(gethistory.NewOptions(msgRepo, gethistory.WithLg(lg)))
	if err != nil {
		return nil, fmt.Errorf("init messages repo: %v", err)
	}

	// initialization v1 handlers
	v1Handlers, err := clientv1.NewHandlers(clientv1.NewOptions(getHistoryUC, clientv1.WithLg(lg)))
	if err != nil {
		return nil, fmt.Errorf("create v1 handlers: %v", err)
	}

	// initialization server
	srv, err := serverclient.New(
		serverclient.NewOptions(
			lg,
			cfg.Servers.Client.Addr,
			cfg.Servers.Client.AllowsOrigins,
			v1Swagger,
			v1Handlers,
			keycloakClient,
			cfg.Servers.Client.RequiredAccess.Resource,
			cfg.Servers.Client.RequiredAccess.Role,
			errHandler.Handle,
		))
	if err != nil {
		return nil, fmt.Errorf("init server: %v", err)
	}

	return srv, nil
}
