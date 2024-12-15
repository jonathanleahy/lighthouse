package service

import (
	"context"
	"errors"

	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/network"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DisputeService DisputeService
	HealthService  HealthService
	GraphqlService GraphqlService
}

func NewResolver(http network.HttpService) *Resolver {
	resolver := &Resolver{
		DisputeService: NewDisputeService(http),
		HealthService:  NewHealthService(),
	}

	resolver.GraphqlService = NewGraphqlService(resolver)
	resolver.GraphqlService.GraphServer.SetErrorPresenter(CustomFormatter)
	return resolver
}

func CustomFormatter(ctx context.Context, err error) *gqlerror.Error {
	var gqlErr *gqlerror.Error
	ok := errors.As(err, &gqlErr)
	if !ok {
		return gqlerror.Errorf("internal server error")
	}

	if gqlErr.Extensions != nil {
		if code, exists := gqlErr.Extensions["code"]; exists && code == "GRAPHQL_VALIDATION_FAILED" {
			gqlErr.Message = "invalid request data"
		}
	}
	return gqlErr
}

