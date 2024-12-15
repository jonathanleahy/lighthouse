package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/logger"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/message"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/utils"
	"net/http"
	"time"
)

var (
	DefaultTimeoutConfig = middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: message.ErrorRequestTimeout,
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			logger.Error(
				message.RequestTimeout,
				utils.GetCIDFromEchoContext(c),
				utils.GetTenantFromEchoContext(c),
				logger.Fields{
					"code": http.StatusGatewayTimeout,
					"uri":  c.Request().RequestURI,
				})
		},
		Timeout: time.Second * time.Duration(env.GetEnvWithDefaultAsInt(env.HttpTimeout, env.DefaultHttpTimeout)),
	}
)

// ConfigTimeout middleware adds a `timeout`
func ConfigTimeout() echo.MiddlewareFunc {
	return middleware.TimeoutWithConfig(DefaultTimeoutConfig)
}

