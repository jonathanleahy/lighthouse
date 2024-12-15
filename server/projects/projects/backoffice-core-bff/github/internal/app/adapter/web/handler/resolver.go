package handler

import "github.com/pismo/backoffice-core-bff/internal/app/domain/service"

var resolver *service.Resolver

func SetResolver(r *service.Resolver) {
	resolver = r
}

