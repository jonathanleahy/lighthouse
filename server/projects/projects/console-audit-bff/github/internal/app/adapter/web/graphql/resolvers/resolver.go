package resolvers

import "github.com/pismo/console-audit-bff/internal/app/domain/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Manager *service.Manager
}

// SearchAudits godoc
// @Summary Return audit by parameters
// @Description Return audit by parameters
// @Tags audit
// @Produce json
// @Router /v1/audit [get]
// @success 200 {object} entity.SearchAudits

