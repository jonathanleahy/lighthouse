package server

import (
	"github.com/labstack/echo/v4"
	"github.com/pismo/backoffice-core-bff/internal/app/adapter/web/handler"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/network"
)

func NewRouter(e *echo.Echo, net network.HttpService) {
	// Health
	e.GET("/health", handler.NewHealthHandler().Get)
	// Graphql
	e.POST("/query", handler.NewGraphqlHandler().Query)
	e.GET("/playground", handler.NewGraphqlHandler().Playground)
}

