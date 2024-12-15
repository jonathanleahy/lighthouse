package server

import (
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/logger"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/network"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/server/middleware"
)

// @title Pismo backoffice-core-bff API
// @version 1.0
// @description This is a backoffice-core-bff api documentation.

// Start Config Router @host localhost:8080
// @BasePath /
func Start(http network.HttpService) chan error {
	e := configEcho()
	NewRouter(e, http)

	errChan := make(chan error, 1)
	go func() {
		logger.Info("Starting application", "", "", nil)
		errChan <- gracehttp.Serve(e.Server)
	}()

	return errChan
}

func configEcho() *echo.Echo {
	e := echo.New()
	// newRequest middleware adds a `x-cid` header to the header and get Tenant and AccountID.
	e.Use(middleware.ConfigRequest())
	//Timeout middleware is used to timeout
	e.Use(middleware.ConfigTimeout())
	//Timeout middleware is used to timeout
	e.Use(echoMiddleware.Recover())
	// Enable tracing middleware
	e.Use(middleware.ConfigOpentelemetry())

	e.Validator = middleware.ConfigValidator()

	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	e.Server.Addr = fmt.Sprintf("%s:%s",
		env.GetEnvWithDefaultAsString(env.Host, env.DefaultHost),
		env.GetEnvWithDefaultAsString(env.Port, env.DefaultPort))

	return e
}

