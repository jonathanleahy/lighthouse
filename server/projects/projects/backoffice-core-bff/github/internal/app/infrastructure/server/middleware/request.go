package middleware

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/request"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/utils"
	"net/http"
	"strings"
)

func ConfigRequest() echo.MiddlewareFunc {
	return newRequest
}

// CidHeader middleware adds a `cid` header to the response.
func newRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cid := c.Request().Header.Get(env.HeaderXCid)
		if cid == "" {
			cid = uuid.New().String()
			c.Request().Header.Set(env.HeaderXCid, cid)
		}

		org := c.Request().Header.Get(env.HeaderXTenant)
		accountID := c.Request().Header.Get(env.HeaderXAccountID)

		rctx := request.RequestContext{
			Tenant:        org,
			AccountID:     accountID,
			Cid:           cid,
			AuditRequest:  utils.GetRequestBody(c),
			Roles:         extractRoles(c.Request().Header),
			CustomHeaders: extractCustomHeaders(c.Request().Header),
		}

		ctx := context.WithValue(c.Request().Context(), env.RequestContext, rctx)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

func extractRoles(header http.Header) []string {
	roleHeaderStr := header.Get(request.HeaderXRoles)
	rolesHeader := strings.Split(roleHeaderStr, request.HeaderXRolesSeparator)
	roles := make([]string, 0)

	for _, role := range rolesHeader {
		roles = append(roles, strings.TrimSpace(role))
	}

	return roles
}

func extractCustomHeaders(headers http.Header) map[string]string {
	extractedHeaders := make(map[string]string)
	for _, header := range request.CustomHeaders {
		extractedHeaders[header] = headers.Get(header)
	}
	return extractedHeaders
}

