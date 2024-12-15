package web

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/pismo/console-audit-bff/internal/app/adapter/web/graphql/entity"

	"github.com/golang/mock/gomock"
	"github.com/pismo/console-audit-bff/internal/app/adapter/web/handler"
	"github.com/pismo/console-audit-bff/internal/app/domain/service"
	"github.com/pismo/psm-sdk/psm"
	"github.com/pismo/psm-sdk/psm/server"
	mockserver "github.com/pismo/psm-sdk/psm/server/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigureRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedServer := mockserver.NewMockServer(ctrl)
	manager := service.NewManager(psm.HttpClient())

	expectedRoutes := []*server.GraphQLRoute{
		{
			Path:                "/query",
			QueryPermissions:    GetQueryPermissions(),
			MutationPermissions: GetMutationPermissions(),
		},
		{
			Path:                "/audit/query",
			QueryPermissions:    GetQueryPermissions(),
			MutationPermissions: GetMutationPermissions(),
		},
	}

	calls := make([]*gomock.Call, 0)

	for i := range expectedRoutes {
		expectedRoute := expectedRoutes[i]
		expectedPermissions := expectedRoute.QueryPermissions
		calls = append(calls, mockedServer.EXPECT().AddRoute(gomock.Any()).AnyTimes())
		calls = append(calls, mockedServer.EXPECT().AddGraphQLRoute(gomock.Any()).Do(func(currentRoute *server.GraphQLRoute) {
			require.NotNil(t, currentRoute.Path)
			assert.Equal(t, expectedPermissions, currentRoute.QueryPermissions)
		}).AnyTimes())
	}
	gomock.InOrder(calls...)
	ConfigureRoutes(mockedServer, manager)
}

func TestConfigureHttpRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedServer := mockserver.NewMockServer(ctrl)
	manager := service.NewManager(psm.HttpClient())

	expectedRoutes := []*server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/whoami",
			Handler: handler.Whoami,
		},
	}

	calls := make([]*gomock.Call, 0)
	for i := range expectedRoutes {
		expectedRoute := expectedRoutes[i]
		expectedHandler := expectedRoute.Handler
		expectedRoute.Handler = nil
		calls = append(calls, mockedServer.EXPECT().AddRoute(gomock.Any()).Do(func(currentRoute *server.Route) {
			require.NotNil(t, currentRoute.Handler)

			currentHandler := currentRoute.Handler
			currentRoute.Handler = nil

			assert.Equal(t, expectedRoute, currentRoute)
			assert.Equal(t, reflect.ValueOf(expectedHandler).Pointer(), reflect.ValueOf(currentHandler).Pointer())
		}))
		calls = append(calls, mockedServer.EXPECT().AddGraphQLRoute(gomock.Any()).AnyTimes())
	}
	gomock.InOrder(calls...)
	ConfigureRoutes(mockedServer, manager)
}

func TestAuditCreateRole(t *testing.T) {
	input := entity.RoleInput{Name: "role-test"}
	req := make(map[string]interface{})
	req["role"] = input
	audit := AuditRoleCreate(req, nil)

	assert.Equal(t, "create", audit.Action())
	assert.Equal(t, "role", audit.Domain())
	assert.Equal(t, input.Name, audit.ID())
}

func TestAuditCreateUpdate(t *testing.T) {
	req := make(map[string]interface{})
	req["id"] = 1
	audit := AuditRoleUpdate(req, nil)

	assert.Equal(t, "update", audit.Action())
	assert.Equal(t, "role", audit.Domain())
	assert.Equal(t, 1, audit.ID())
}

func TestAuditCreateDelete(t *testing.T) {
	req := make(map[string]interface{})
	req["id"] = 1
	audit := AuditRoleDelete(req, nil)

	assert.Equal(t, "delete", audit.Action())
	assert.Equal(t, "role", audit.Domain())
	assert.Equal(t, 1, audit.ID())
}

func TestAuditCreateFeature(t *testing.T) {
	input := entity.FeatureInput{Name: "role-test"}
	req := make(map[string]interface{})
	req["feature"] = input
	audit := AuditFeatureCreate(req, nil)

	assert.Equal(t, "create", audit.Action())
	assert.Equal(t, "feature", audit.Domain())
	assert.Equal(t, input.Name, audit.ID())
}

func TestAuditUpdateFeature(t *testing.T) {
	req := make(map[string]interface{})
	req["id"] = 1
	audit := AuditFeatureUpdate(req, nil)

	assert.Equal(t, "update", audit.Action())
	assert.Equal(t, "feature", audit.Domain())
	assert.Equal(t, 1, audit.ID())
}

func TestAuditDeleteFeature(t *testing.T) {
	req := make(map[string]interface{})
	req["id"] = 1
	audit := AuditFeatureDelete(req, nil)

	assert.Equal(t, "delete", audit.Action())
	assert.Equal(t, "feature", audit.Domain())
	assert.Equal(t, 1, audit.ID())
}

func TestAuditRoleAttachUser(t *testing.T) {
	req := make(map[string]interface{})
	req["id"] = 1
	req["email"] = "pismo@pismo.io"
	audit := AuditRoleAttachUser(req, nil)

	assert.Equal(t, "create", audit.Action())
	assert.Equal(t, "role-attach-to-user", audit.Domain())
	assert.Equal(t, fmt.Sprintf("%s-%s", req["id"], req["email"]), audit.ID())
}

func TestAuditRoleDetachUser(t *testing.T) {
	req := make(map[string]interface{})
	req["id"] = 1
	req["email"] = "pismo@pismo.io"
	audit := AuditRoleDetachUser(req, nil)

	assert.Equal(t, "create", audit.Action())
	assert.Equal(t, "role-detach-to-user", audit.Domain())
	assert.Equal(t, fmt.Sprintf("%s-%s", req["id"], req["email"]), audit.ID())
}

func TestAuditRoleAttachFeature(t *testing.T) {
	req := make(map[string]interface{})
	req["id"] = 1
	req["featureID"] = 1
	audit := AuditRoleAttachFeature(req, nil)

	assert.Equal(t, "create", audit.Action())
	assert.Equal(t, "role-attach-to-feature", audit.Domain())
	assert.Equal(t, fmt.Sprintf("%s-%s", req["id"], req["featureID"]), audit.ID())
}

func TestAuditRoleDetachFeature(t *testing.T) {
	req := make(map[string]interface{})
	req["id"] = 1
	req["featureID"] = 1
	audit := AuditRoleDetachFeature(req, nil)

	assert.Equal(t, "delete", audit.Action())
	assert.Equal(t, "role-detach-to-feature", audit.Domain())
	assert.Equal(t, fmt.Sprintf("%s-%s", req["id"], req["featureID"]), audit.ID())
}

func TestAuditCreateEndpoint(t *testing.T) {
	input := entity.EndpointInput{ServiceName: "test-api", Method: "POST", Path: "/users/"}
	req := make(map[string]interface{})
	req["endpoint"] = input
	audit := AuditEndpointCreate(req, nil)
	data := fmt.Sprintf("%s-%s-%s", input.ServiceName, input.Path, input.Method)

	assert.Equal(t, "create", audit.Action())
	assert.Equal(t, "endpoint", audit.Domain())
	assert.Equal(t, data, audit.ID())
}

func TestAuditUpdateEndpoint(t *testing.T) {
	req := make(map[string]interface{})
	req["id"] = 1
	audit := AuditEndpointUpdate(req, nil)

	assert.Equal(t, "update", audit.Action())
	assert.Equal(t, "endpoint", audit.Domain())
	assert.Equal(t, 1, audit.ID())
}

func TestAuditDeleteEndpoint(t *testing.T) {
	req := make(map[string]interface{})
	req["id"] = 1
	audit := AuditEndpointDelete(req, nil)

	assert.Equal(t, "delete", audit.Action())
	assert.Equal(t, "endpoint", audit.Domain())
	assert.Equal(t, 1, audit.ID())
}

