package web

import (
	"fmt"
	"net/http"

	"github.com/pismo/console-audit-bff/internal/app/adapter/web/graphql/entity"
	"github.com/pismo/psm-sdk/psm/server/audit"

	"github.com/pismo/console-audit-bff/internal/app/adapter/web/graphql/generated"
	"github.com/pismo/console-audit-bff/internal/app/adapter/web/graphql/resolvers"
	"github.com/pismo/console-audit-bff/internal/app/adapter/web/handler"
	"github.com/pismo/console-audit-bff/internal/app/domain/service"
	"github.com/pismo/psm-sdk/psm/server"
)

// ConfigureRoutes godoc
// @Title Pismo Audit-BFF
// @Summary Return graphql data
// @Description Return graphql data
// @Tags graphql
// @Produce json
// @success 200
// @success 500
// @Param "x-roles" header string true "Roles"
// @Param "x-tenant" header string true "Tenant"
// @Param query body string true "query"
// @Router /query [post]
func ConfigureRoutes(srv server.Server, manager *service.Manager) {
	handler.SetManager(manager)

	srv.AddRoute(&server.Route{
		Method:  http.MethodGet,
		Path:    "/whoami",
		Handler: handler.Whoami,
	})

	srv.AddGraphQLRoute(&server.GraphQLRoute{
		Path: "/query",
		Schema: generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &resolvers.Resolver{
					Manager: manager,
				},
			},
		),
		QueryPermissions:    GetQueryPermissions(),
		MutationPermissions: GetMutationPermissions(),
	})

	srv.AddGraphQLRoute(&server.GraphQLRoute{
		Path: "/audit/query",
		Schema: generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &resolvers.Resolver{
					Manager: manager,
				},
			},
		),
		QueryPermissions:    GetQueryPermissions(),
		MutationPermissions: GetMutationPermissions(),
	})
}

func GetQueryPermissions() []*server.GraphQLPermission {
	return []*server.GraphQLPermission{
		{
			Operation: "list_roles",
			Roles: []string{
				"controlc-admin",
				"controlc-user-manager",
				"controlc-auditor",
				"controlc-setup-advanced",
				"controlc-setup-advanced-viewer",
				"controlc-setup-operator",
				"controlc-setup-viewer",
				"controlc-accounts-operator",
				"controlc-accounts-viewer",
				"controlc-backoffice-operator",
				"controlc-backoffice-viewer",
				"controlc-seller-mgmt-admin",
			},
		},
		{
			Operation: "find_audit_by_id",
			Roles:     []string{"controlc-admin", "controlc-auditor"},
		},
		{
			Operation: "search_audit",
			Roles:     []string{"controlc-admin", "controlc-auditor"},
		},
		{
			Operation: "list_user_roles",
			Roles: []string{
				"controlc-admin",
				"controlc-user-manager",
				"controlc-auditor",
				"controlc-setup-advanced",
				"controlc-setup-advanced-viewer",
				"controlc-setup-operator",
				"controlc-setup-viewer",
				"controlc-accounts-operator",
				"controlc-accounts-viewer",
				"controlc-backoffice-operator",
				"controlc-backoffice-viewer",
				"controlc-seller-mgmt-admin",
			},
		},
		{
			Operation: "search_roles",
			Roles: []string{
				"controlc-admin",
				"controlc-user-manager",
				"controlc-auditor",
				"controlc-setup-advanced",
				"controlc-setup-advanced-viewer",
				"controlc-setup-operator",
				"controlc-setup-viewer",
				"controlc-accounts-operator",
				"controlc-accounts-viewer",
				"controlc-backoffice-operator",
				"controlc-backoffice-viewer",
				"controlc-seller-mgmt-admin",
			},
		},
		{
			Operation: "find_feature_by_id",
			Roles:     []string{"owner"},
		},
		{
			Operation: "search_feature",
			Roles:     []string{"owner"},
		},
		{
			Operation: "find_endpoint_by_id",
			Roles:     []string{"owner"},
		},
		{
			Operation: "list_endpoint",
			Roles:     []string{"owner"},
		},
	}
}

func GetMutationPermissions() []*server.GraphQLPermission {
	return []*server.GraphQLPermission{
		{
			Operation: "create_role",
			Roles:     []string{"owner"},
			Audit:     AuditRoleCreate,
		},
		{
			Operation: "update_role",
			Roles:     []string{"owner"},
			Audit:     AuditRoleUpdate,
		},
		{
			Operation: "delete_role",
			Roles:     []string{"owner"},
			Audit:     AuditRoleDelete,
		},
		{
			Operation: "create_feature",
			Roles:     []string{"owner"},
			Audit:     AuditFeatureCreate,
		},
		{
			Operation: "update_feature",
			Roles:     []string{"owner"},
			Audit:     AuditFeatureUpdate,
		},
		{
			Operation: "delete_feature",
			Roles:     []string{"owner"},
			Audit:     AuditFeatureDelete,
		},
		{
			Operation: "attach_role_to_user",
			Roles:     []string{"owner"},
			Audit:     AuditRoleAttachUser,
		},
		{
			Operation: "detach_role_to_user",
			Roles:     []string{"owner"},
			Audit:     AuditRoleDetachUser,
		},
		{
			Operation: "attach_role_to_feature",
			Roles:     []string{"owner"},
			Audit:     AuditRoleAttachFeature,
		},
		{
			Operation: "detach_role_to_feature",
			Roles:     []string{"owner"},
			Audit:     AuditRoleDetachFeature,
		},
		{
			Operation: "create_endpoint",
			Roles:     []string{"owner"},
			Audit:     AuditEndpointCreate,
		},
		{
			Operation: "update_endpoint",
			Roles:     []string{"owner"},
			Audit:     AuditEndpointUpdate,
		},
		{
			Operation: "delete_endpoint",
			Roles:     []string{"owner"},
			Audit:     AuditEndpointDelete,
		},
	}
}

func AuditRoleCreate(req map[string]interface{}, res interface{}) *audit.Audit {
	return audit.Create("role", nil, req["role"].(entity.RoleInput).Name)
}

func AuditRoleUpdate(req map[string]interface{}, res interface{}) *audit.Audit {
	return audit.Update("role", nil, req["id"])
}

func AuditRoleDelete(req map[string]interface{}, res interface{}) *audit.Audit {
	return audit.Delete("role", nil, req["id"])
}

func AuditFeatureCreate(req map[string]interface{}, res interface{}) *audit.Audit {
	return audit.Create("feature", nil, req["feature"].(entity.FeatureInput).Name)
}

func AuditFeatureUpdate(req map[string]interface{}, res interface{}) *audit.Audit {
	return audit.Update("feature", nil, req["id"])
}

func AuditFeatureDelete(req map[string]interface{}, res interface{}) *audit.Audit {
	return audit.Delete("feature", nil, req["id"])
}

func AuditRoleAttachUser(req map[string]interface{}, res interface{}) *audit.Audit {
	return audit.Create("role-attach-to-user", nil, fmt.Sprintf("%s-%s", req["id"], req["email"]))
}

func AuditRoleDetachUser(req map[string]interface{}, res interface{}) *audit.Audit {
	return audit.Create("role-detach-to-user", nil, fmt.Sprintf("%s-%s", req["id"], req["email"]))
}

func AuditRoleAttachFeature(req map[string]interface{}, res interface{}) *audit.Audit {
	return audit.Create("role-attach-to-feature", nil, fmt.Sprintf("%s-%s", req["id"], req["featureID"]))
}

func AuditRoleDetachFeature(req map[string]interface{}, res interface{}) *audit.Audit {
	return audit.Delete("role-detach-to-feature", nil, fmt.Sprintf("%s-%s", req["id"], req["featureID"]))
}

func AuditEndpointCreate(req map[string]interface{}, res interface{}) *audit.Audit {
	data := fmt.Sprintf("%s-%s-%s", req["endpoint"].(entity.EndpointInput).ServiceName, req["endpoint"].(entity.EndpointInput).Path, req["endpoint"].(entity.EndpointInput).Method)
	return audit.Create("endpoint", nil, data)
}

func AuditEndpointUpdate(req map[string]interface{}, res interface{}) *audit.Audit {
	return audit.Update("endpoint", nil, req["id"])
}

func AuditEndpointDelete(req map[string]interface{}, res interface{}) *audit.Audit {
	return audit.Delete("endpoint", nil, req["id"])
}

