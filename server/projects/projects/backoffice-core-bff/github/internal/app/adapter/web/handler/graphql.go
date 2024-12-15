package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/apierror"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/message"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/utils"
	"io"
	"net/http"
	"strings"
)

const (
	schemaQuery = "__schema"
)

type HTTPError struct {
	Code     int         `json:"-"`
	Message  interface{} `json:"message"`
	Internal error       `json:"-"` // Stores the error returned by an external dependency
}

type GraphqlHandler interface {
	Query(c echo.Context) error
	Playground(c echo.Context) error
}

type graphqlHandler struct{}

func NewGraphqlHandler() GraphqlHandler {
	return &graphqlHandler{}
}

// Query godoc
// @Summary Return graphql data
// @Description Return graphql data
// @Tags graphql
// @Produce json
// @success 200
// @success 500
// @Param "x-roles" header string true "Roles"
// @Param "x-tenancy" header string true "Tenancy"
// @Param query body string true "query"
// @Router /query [post]
func (h graphqlHandler) Query(c echo.Context) error {
	err := verifyQueryIsAllowed(c)
	if err != nil {
		apiError := err.(apierror.ApiError)
		httpError := HTTPError{
			Message: apiError.GetMessage().UserMessage,
		}
		return c.JSON(http.StatusForbidden, httpError)
	}
	resolver.GraphqlService.GraphServer.ServeHTTP(c.Response(), c.Request())
	return nil
}

// Playground godoc
// @Summary Return graphql playground
// @Description Return graphql playground
// @Tags graphql
// @Router /playground [get]
func (h graphqlHandler) Playground(c echo.Context) error {
	resolver.GraphqlService.PlaygroundService.ServeHTTP(c.Response(), c.Request())
	return nil
}

func verifyQueryIsAllowed(c echo.Context) error {
	if env.GetEnvWithDefaultAsString(env.Env, env.DefaultEnv) == "dev" {
		return nil
	}

	byteBody, _ := io.ReadAll(c.Request().Body)
	var response = make(map[string]interface{})
	err := json.Unmarshal(byteBody, &response)
	if err != nil {
		return err
	}

	c.Request().Body = io.NopCloser(bytes.NewBuffer(byteBody))

	if s, ok := response["query"].(string); ok {
		if strings.Contains(s, schemaQuery) {
			errorMessage := message.ErrQuery
			ctx := c.Request().Context()
			return apierror.NewError(ctx, errors.New(errorMessage.UserMessage), errorMessage,
				"api_error", utils.GetCIDFromContext(ctx), utils.GetTenantFromContext(ctx), http.StatusForbidden)
		}
	}
	return nil
}

