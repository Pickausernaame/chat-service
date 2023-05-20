package main

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	keycloakclient "github.com/Pickausernaame/chat-service/internal/clients/keycloak"
	"github.com/Pickausernaame/chat-service/internal/config"
	chatsrepo "github.com/Pickausernaame/chat-service/internal/repositories/chats"
	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	problemsrepo "github.com/Pickausernaame/chat-service/internal/repositories/problems"
	"github.com/Pickausernaame/chat-service/internal/server"
	clientv1 "github.com/Pickausernaame/chat-service/internal/server-client/v1"
	"github.com/Pickausernaame/chat-service/internal/server/errhandler"
	"github.com/Pickausernaame/chat-service/internal/services/outbox"
	gethistory "github.com/Pickausernaame/chat-service/internal/usecases/client/get-history"
	sendmessage "github.com/Pickausernaame/chat-service/internal/usecases/client/send-message"
)

const nameServerClient = "server-client"

type Transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

func initServerClient(
	cfg config.Config,
	keycloakClient *keycloakclient.Client,
	msgRepo *messagesrepo.Repo,
	chatRepo *chatsrepo.Repo,
	problemRepo *problemsrepo.Repo,
	txtr Transactor,
	outbox *outbox.Service,
) (*server.Server, error) {
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
	getHistoryUC, err := gethistory.New(gethistory.NewOptions(msgRepo))
	if err != nil {
		return nil, fmt.Errorf("init getHistory usecase: %v", err)
	}

	// initialization sendMessage useCase
	sendMessageUC, err := sendmessage.New(sendmessage.NewOptions(chatRepo, msgRepo, outbox, problemRepo, txtr))
	if err != nil {
		return nil, fmt.Errorf("init sendMessage usecase: %v", err)
	}

	// initialization v1 handlers
	v1Handlers, err := clientv1.NewHandlers(clientv1.NewOptions(getHistoryUC, sendMessageUC))
	if err != nil {
		return nil, fmt.Errorf("create v1 handlers: %v", err)
	}

	// initialization server
	srv, err := server.New(
		server.NewOptions(
			lg,
			cfg.Servers.Client.Addr,
			cfg.Servers.Client.AllowsOrigins,
			v1Swagger,
			func(router *echo.Group) { clientv1.RegisterHandlers(router, v1Handlers) },
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
