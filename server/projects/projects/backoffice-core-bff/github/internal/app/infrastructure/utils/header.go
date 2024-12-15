package utils

import (
	"bytes"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/request"
	"io"
	"strings"
)

func GetRequestIP(xFowardedFor string) string {
	ss := strings.Split(xFowardedFor, ",")
	return ss[0]
}

func GetRequestBody(c echo.Context) string {
	// Read the Body content
	var bodyBytes []byte
	if c.Request().Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request().Body)
	}

	// Restore the io.ReadCloser to its original state
	c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return string(bodyBytes)
}

func GetCIDFromEchoContext(c echo.Context) string {
	return c.Request().Header.Get(env.HeaderXCid)
}

func GetTenantFromEchoContext(c echo.Context) string {
	return c.Request().Header.Get(env.HeaderXTenant)
}

func GetCIDFromContext(ctx context.Context) string {
	return ctx.Value(env.RequestContext).(request.RequestContext).Cid
}

func GetTenantFromContext(ctx context.Context) string {
	return ctx.Value(env.RequestContext).(request.RequestContext).Tenant
}

