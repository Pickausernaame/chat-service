package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	keycloakclient "github.com/Pickausernaame/chat-service/internal/clients/keycloak"
	"github.com/Pickausernaame/chat-service/internal/config"
	"github.com/Pickausernaame/chat-service/internal/server"
	managerv1 "github.com/Pickausernaame/chat-service/internal/server-manager/v1"
	"github.com/Pickausernaame/chat-service/internal/server/errhandler"
	manager_load "github.com/Pickausernaame/chat-service/internal/services/manager-load"
	managerpool "github.com/Pickausernaame/chat-service/internal/services/manager-pool"
	canreceiveproblems "github.com/Pickausernaame/chat-service/internal/usecases/manager/can-receive-problems"
	setreadyreceiveproblems "github.com/Pickausernaame/chat-service/internal/usecases/manager/set-ready-receive-problems"
)

const nameServerManager = "server-manager"

func initServerManager(
	cfg config.Config,
	keycloakClient *keycloakclient.Client,
	managerLoadService *manager_load.Service,
	managerPool managerpool.Pool,
) (*server.Server, error) {
	lg := zap.L().Named(nameServerManager)

	// getting specification
	v1Swagger, err := managerv1.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("getting swagger: %v", err)
	}

	// initialization errorHandler
	errHandler, err := errhandler.New(errhandler.NewOptions(lg, cfg.Global.IsProd(), errhandler.ResponseBuilder))
	if err != nil {
		return nil, fmt.Errorf("init errror handler: %v", err)
	}

	// creating useCases
	// initialization canReciveProblems useCase
	canReciveProblems, err := canreceiveproblems.New(canreceiveproblems.NewOptions(managerLoadService, managerPool))
	if err != nil {
		return nil, fmt.Errorf("init canReciveProblems usecase: %v", err)
	}

	// initialization setReadyReceiveProblems useCase
	setReadyReceiveProblems, err := setreadyreceiveproblems.New(setreadyreceiveproblems.NewOptions(managerLoadService, managerPool))
	if err != nil {
		return nil, fmt.Errorf("init setReadyReceiveProblems usecase: %v", err)
	}

	// initialization v1 handlers
	v1Handlers, err := managerv1.NewHandlers(managerv1.NewOptions(canReciveProblems, setReadyReceiveProblems))
	if err != nil {
		return nil, fmt.Errorf("create v1 handlers: %v", err)
	}

	// initialization server
	srv, err := server.New(
		server.NewOptions(
			lg,
			cfg.Servers.Manager.Addr,
			cfg.Servers.Manager.AllowsOrigins,
			v1Swagger,
			func(router *echo.Group) {
				managerv1.RegisterHandlers(router, v1Handlers)
			},
			keycloakClient,
			cfg.Servers.Manager.RequiredAccess.Resource,
			cfg.Servers.Manager.RequiredAccess.Role,
			errHandler.Handle,
		))
	if err != nil {
		return nil, fmt.Errorf("init server: %v", err)
	}

	return srv, nil
}
