package service

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/pismo/backoffice-core-bff/internal/app/adapter/graphql/generated"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/apierror"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/request"
	"net/http"
)

type GraphqlService struct {
	GraphServer       *handler.Server
	PlaygroundService http.HandlerFunc
}

func NewGraphqlService(resolver *Resolver) GraphqlService {
	graphServer := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers:  resolver,
				Directives: generated.DirectiveRoot{HasAnyRole: HasAnyRole},
			},
		),
	)

	return GraphqlService{
		GraphServer:       graphServer,
		PlaygroundService: playground.Handler("GraphQL", "/query"),
	}
}

func HasAnyRole(ctx context.Context, obj interface{}, next graphql.Resolver, roles []string) (res interface{}, err error) {
	rctx := ctx.Value(env.RequestContext).(request.RequestContext)
	if rctx.HasAnyRole(roles...) {
		return next(ctx)
	}

	return nil, apierror.ForbiddenError(ctx)
}

