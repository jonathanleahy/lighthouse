package app

import (
	"context"
	"github.com/pismo/console-audit-bff/internal/app/adapter/web"
	"github.com/pismo/console-audit-bff/internal/app/domain/service"
	"github.com/pismo/psm-sdk/psm/schema"
	"github.com/pismo/psm-sdk/psm/server"
)

type (
	Application interface {
		Start() error
		Stop()
	}

	application struct {
		server  server.Server
		manager *service.Manager
	}
)

func NewApplication(httpClient schema.HttpClient) (Application, error) {
	manager := service.NewManager(httpClient)
	srv := configureServer(manager)

	return &application{
		server:  srv,
		manager: manager,
	}, nil
}

func configureServer(manager *service.Manager) server.Server {
	srv := server.NewDefaultServer()
	web.ConfigureRoutes(srv, manager)
	return srv
}

func (a *application) Start() error {
	return a.server.Start(context.Background())
}

func (a *application) Stop() {
	_ = a.server.Stop(context.Background())
}

