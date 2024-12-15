package resolvers

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/pismo/console-audit-bff/internal/app/adapter/web/graphql/entity"

	"github.com/golang/mock/gomock"
	"github.com/pismo/console-audit-bff/internal/app/adapter/web/presenter"
	"github.com/pismo/console-audit-bff/internal/app/domain/service"
	"github.com/pismo/console-audit-bff/internal/app/domain/service/mock"
	"github.com/pismo/psm-sdk/psm"
	"github.com/pismo/psm-sdk/psm/factory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryResolver_FindAuditByID_Success(t *testing.T) {
	ctx := context.Background()

	device := pointer("Mozilla")
	deviceIP := pointer("179.177.200.229")

	lat := pointer(-23.592049)
	long := pointer(-46.685087)

	res := &presenter.Audit{
		ID: 1,
		Operation: &presenter.Operation{
			Tenant:   "TN-PISMO",
			Action:   "create",
			Domain:   "user",
			DomainID: "crm-operator@pismo.io",
			CID:      "CID-crm-operator",
			Date:     "2022-01-01 12:00:00",
		},
		User: &presenter.User{
			Email: "pismo@pismo.com",
			Roles: []string{"crm-operator"},
		},
		UserAgent: &presenter.UserAgent{
			Device:   device,
			DeviceIp: deviceIP,
		},
		Localization: &presenter.Localization{
			Latitude:  lat,
			Longitude: long,
		},
		Http: &presenter.Http{
			Code:     200,
			Request:  `{"name": "crm-operator", "email": "crm-operator@pismo.io", "roles": ["crm-operator"]}`,
			Response: `{"id": 1}`,
		},
	}

	psm.Configure("console-audit-bff", factory.Default())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedService := mock.NewMockAuditService(ctrl)
	mockedService.EXPECT().
		FindAuditByID(ctx, gomock.Any()).
		Return(res, nil)

	r := &queryResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}

	// test
	response, err := r.FindAuditByID(ctx, 1)
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, 1, response.ID)
	require.Equal(t, "TN-PISMO", response.Operation.Tenant)
	require.Equal(t, "create", response.Operation.Action)
	require.Equal(t, "user", response.Operation.Domain)
	require.Equal(t, "crm-operator@pismo.io", response.Operation.DomainID)
	require.Equal(t, "2022-01-01 12:00:00", response.Operation.Date)
	require.Equal(t, "pismo@pismo.com", response.User.Email)
	require.Equal(t, "crm-operator", response.User.Roles[0])
	require.Equal(t, device, response.UserAgent.Device)
	require.Equal(t, deviceIP, response.UserAgent.DeviceIP)
	require.Equal(t, lat, response.Localization.Latitude)
	require.Equal(t, long, response.Localization.Longitude)
	require.Equal(t, 200, response.HTTP.Code)
	require.Equal(t, `{"name": "crm-operator", "email": "crm-operator@pismo.io", "roles": ["crm-operator"]}`, response.HTTP.Request)
	require.Equal(t, `{"id": 1}`, response.HTTP.Response)
}

func TestQueryResolver_FindAuditByID_Error(t *testing.T) {
	ctx := context.Background()

	psm.Configure("console-audit-bff", factory.Default())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedService := mock.NewMockAuditService(ctrl)
	mockedService.EXPECT().
		FindAuditByID(ctx, gomock.Any()).
		Return(nil, errors.New("err"))

	r := &queryResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}

	response, err := r.FindAuditByID(ctx, 1)
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, "err", err.Error())
}

func TestQueryResolver_SearchAudit_Success(t *testing.T) {
	ctx := context.Background()

	device := pointer("Mozilla")
	deviceIP := pointer("179.177.200.229")

	lat := pointer(-23.592049)
	long := pointer(-46.685087)

	res := &presenter.SearchAudit{
		CurrentPage: 1,
		Pages:       1,
		PerPage:     1,
		TotalItems:  1,
		Items: []*presenter.Audit{
			{
				ID: 1,
				Operation: &presenter.Operation{
					Tenant:   "TN-PISMO",
					Action:   "create",
					Domain:   "user",
					DomainID: "crm-operator@pismo.io",
					CID:      "CID-crm-operator",
					Date:     "2022-01-01 12:00:00",
				},
				User: &presenter.User{
					Email: "pismo@pismo.com",
					Roles: []string{"crm-operator"},
				},
				UserAgent: &presenter.UserAgent{
					Device:   device,
					DeviceIp: deviceIP,
				},
				Localization: &presenter.Localization{
					Latitude:  lat,
					Longitude: long,
				},
				Http: &presenter.Http{
					Code:     200,
					Request:  `{"name": "crm-operator", "email": "crm-operator@pismo.io", "roles": ["crm-operator"]}`,
					Response: `{"id": 1}`,
				},
			},
		},
	}

	psm.Configure("console-audit-bff", factory.Default())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedService := mock.NewMockAuditService(ctrl)
	mockedService.EXPECT().
		SearchAudit(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(res, nil)

	r := &queryResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}

	// test
	response, err := r.SearchAudit(ctx, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, 1, response.CurrentPage)
	require.Equal(t, 1, response.Pages)
	require.Equal(t, 1, response.PerPage)
	require.Equal(t, 1, response.TotalItems)
	require.Equal(t, 1, response.Audits[0].ID)
	require.Equal(t, "TN-PISMO", response.Audits[0].Operation.Tenant)
	require.Equal(t, "create", response.Audits[0].Operation.Action)
	require.Equal(t, "user", response.Audits[0].Operation.Domain)
	require.Equal(t, "crm-operator@pismo.io", response.Audits[0].Operation.DomainID)
	require.Equal(t, "2022-01-01 12:00:00", response.Audits[0].Operation.Date)
	require.Equal(t, "pismo@pismo.com", response.Audits[0].User.Email)
	require.Equal(t, "crm-operator", response.Audits[0].User.Roles[0])
	require.Equal(t, device, response.Audits[0].UserAgent.Device)
	require.Equal(t, deviceIP, response.Audits[0].UserAgent.DeviceIP)
	require.Equal(t, lat, response.Audits[0].Localization.Latitude)
	require.Equal(t, long, response.Audits[0].Localization.Longitude)
	require.Equal(t, 200, response.Audits[0].HTTP.Code)
	require.Equal(t, `{"name": "crm-operator", "email": "crm-operator@pismo.io", "roles": ["crm-operator"]}`, response.Audits[0].HTTP.Request)
	require.Equal(t, `{"id": 1}`, response.Audits[0].HTTP.Response)
}

func TestQueryResolver_SearchAudit_Error(t *testing.T) {
	ctx := context.Background()

	psm.Configure("console-audit-bff", factory.Default())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedService := mock.NewMockAuditService(ctrl)
	mockedService.EXPECT().
		SearchAudit(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("err"))

	r := &queryResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}

	response, err := r.SearchAudit(ctx, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, "err", err.Error())
}

func TestQueryResolver_ListRoles_Success(t *testing.T) {
	ctx := context.Background()

	psm.Configure("console-audit-bff", factory.Default())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedService := mock.NewMockAuditService(ctrl)
	mockedService.EXPECT().
		ListRoles(ctx).
		Return([]*entity.Role{}, nil)

	r := &queryResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}

	// test
	response, err := r.ListRoles(ctx)
	require.NoError(t, err)
	require.NotNil(t, response)
}

func TestQueryResolver_ListRoles_Error(t *testing.T) {
	ctx := context.Background()

	psm.Configure("console-audit-bff", factory.Default())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedService := mock.NewMockAuditService(ctrl)
	mockedService.EXPECT().
		ListRoles(ctx).
		Return(nil, errors.New("err"))

	r := &queryResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}

	response, err := r.ListRoles(ctx)
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, "err", err.Error())
}

func TestResolver_Query(t *testing.T) {
	resolver := &Resolver{
		Manager: &service.Manager{
			AuditService: nil,
		},
	}

	queryResolver := resolver.Query()
	assert.NotNil(t, queryResolver)
}

func TestResolver_Mutation(t *testing.T) {
	resolver := &Resolver{
		Manager: &service.Manager{
			AuditService: nil,
		},
	}

	mutationResolver := resolver.Mutation()
	assert.NotNil(t, mutationResolver)
}

func pointer[T any](value T) *T {
	return &value
}

func TestQueryResolver_ListUserRoles_Success(t *testing.T) {
	ctx := context.Background()

	res := &presenter.UserRoles{
		Roles: []string{"admin", "crm-operator"},
	}

	psm.Configure("console-audit-bff", factory.Default())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedService := mock.NewMockAuditService(ctrl)
	mockedService.EXPECT().
		ListUserRoles(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(res, nil)

	r := &queryResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}

	// test
	response, err := r.ListUserRoles(ctx, nil, nil, nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, response)
}

func TestQueryResolver_ListUserRoles_Error(t *testing.T) {
	ctx := context.Background()

	psm.Configure("console-audit-bff", factory.Default())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedService := mock.NewMockAuditService(ctrl)
	mockedService.EXPECT().
		ListUserRoles(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("err"))

	r := &queryResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}

	response, err := r.ListUserRoles(ctx, nil, nil, nil, nil, nil)
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, "err", err.Error())
}
func TestQueryResolver_SearchRoles_Success(t *testing.T) {
	ctx := context.Background()

	res := presenter.Roles{
		"organization": map[string]map[string][]string{"access_key": {"read": []string{"admin"}, "view": []string{"admin", "crm-operator"}}},
	}

	psm.Configure("console-audit-bff", factory.Default())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedService := mock.NewMockAuditService(ctrl)
	mockedService.EXPECT().
		SearchRoles(ctx, gomock.Any(), gomock.Any()).
		Return(res, nil)

	r := &queryResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}

	// test
	response, err := r.SearchRoles(ctx, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, response)
}

func TestQueryResolver_SearchRoles_Error(t *testing.T) {
	ctx := context.Background()

	psm.Configure("console-audit-bff", factory.Default())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedService := mock.NewMockAuditService(ctrl)
	mockedService.EXPECT().
		SearchRoles(ctx, gomock.Any(), gomock.Any()).
		Return(nil, errors.New("err"))

	r := &queryResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}

	response, err := r.SearchRoles(ctx, nil, nil)
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, "err", err.Error())
}

func Test_mutationResolver_CreateRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	input := entity.RoleInput{}
	mockedService := mock.NewMockAuditService(ctrl)
	gomock.InOrder(
		mockedService.EXPECT().CreateRole(ctx, gomock.Any()).Return(nil),
		mockedService.EXPECT().CreateRole(ctx, gomock.Any()).Return(errors.New("Duplicated role")),
		mockedService.EXPECT().CreateRole(ctx, gomock.Any()).Return(errors.New("Unknown error")),
	)
	tests := []struct {
		name    string
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "success",
			want:    "Role created successfully",
			wantErr: assert.NoError,
		},
		{
			name:    "bad request",
			want:    "Duplicated role",
			wantErr: assert.NoError,
		},
		{
			name:    "fail",
			want:    "Failed to create role",
			wantErr: assert.Error,
		},
	}
	r := &mutationResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.CreateRole(ctx, input)
			if !tt.wantErr(t, err, fmt.Sprintf("CreateRole(%v, %v)", ctx, input)) {
				return
			}
			assert.Equalf(t, tt.want, got, "CreateRole(%v, %v)", ctx, input)
		})
	}
}

func Test_mutationResolver_UpdateRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	input := entity.RoleInput{ID: new(int)}
	mockedService := mock.NewMockAuditService(ctrl)
	gomock.InOrder(
		mockedService.EXPECT().UpdateRole(ctx, gomock.Any()).Return(nil),
		mockedService.EXPECT().UpdateRole(ctx, gomock.Any()).Return(errors.New("Duplicated role")),
		mockedService.EXPECT().UpdateRole(ctx, gomock.Any()).Return(errors.New("Unknown error")),
	)
	tests := []struct {
		name    string
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "success",
			want:    "Role updated successfully",
			wantErr: assert.NoError,
		},
		{
			name:    "bad request",
			want:    "Duplicated role",
			wantErr: assert.NoError,
		},
		{
			name:    "fail",
			want:    "Failed to update role",
			wantErr: assert.Error,
		},
	}
	r := &mutationResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.UpdateRole(ctx, input)
			if !tt.wantErr(t, err, fmt.Sprintf("UpdateRole")) {
				return
			}
			assert.Equalf(t, tt.want, got, "UpdateRole")
		})
	}
}

func Test_mutationResolver_DeleteRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockedService := mock.NewMockAuditService(ctrl)
	gomock.InOrder(
		mockedService.EXPECT().DeleteRole(ctx, gomock.Any()).Return(nil),
		mockedService.EXPECT().DeleteRole(ctx, gomock.Any()).Return(errors.New("Unknown error")),
	)
	tests := []struct {
		name    string
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "success",
			want:    "Role deleted successfully",
			wantErr: assert.NoError,
		},
		{
			name:    "fail",
			want:    "Failed to delete role",
			wantErr: assert.Error,
		},
	}
	r := &mutationResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.DeleteRole(ctx, 1)
			if !tt.wantErr(t, err, fmt.Sprintf("DeleteRole")) {
				return
			}
			assert.Equalf(t, tt.want, got, "DeleteRole")
		})
	}
}

func Test_mutationResolver_CreateFeature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	input := entity.FeatureInput{}
	mockedService := mock.NewMockAuditService(ctrl)
	gomock.InOrder(
		mockedService.EXPECT().CreateFeature(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
		mockedService.EXPECT().CreateFeature(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Duplicated feature")),
		mockedService.EXPECT().CreateFeature(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Unknown error")),
	)
	tests := []struct {
		name    string
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "success",
			want:    "Feature created successfully",
			wantErr: assert.NoError,
		},
		{
			name:    "bad request",
			want:    "Duplicated feature",
			wantErr: assert.NoError,
		},
		{
			name:    "fail",
			want:    "Failed to create feature",
			wantErr: assert.Error,
		},
	}
	r := &mutationResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.CreateFeature(ctx, input, nil, make([]entity.AuditAction, 0), []string{"domain"}, nil, nil)
			if !tt.wantErr(t, err, fmt.Sprintf("CreateFeature(%v, %v)", ctx, input)) {
				return
			}
			assert.Equalf(t, tt.want, got, "CreateFeature(%v, %v)", ctx, input)
		})
	}
}

func Test_mutationResolver_UpdateFeature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	input := entity.FeatureInput{}
	mockedService := mock.NewMockAuditService(ctrl)
	gomock.InOrder(
		mockedService.EXPECT().UpdateFeature(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
		mockedService.EXPECT().UpdateFeature(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Duplicated feature")),
		mockedService.EXPECT().UpdateFeature(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Unknown error")),
	)
	tests := []struct {
		name    string
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "success",
			want:    "Feature updated successfully",
			wantErr: assert.NoError,
		},
		{
			name:    "bad request",
			want:    "Duplicated feature",
			wantErr: assert.NoError,
		},
		{
			name:    "fail",
			want:    "Failed to update feature",
			wantErr: assert.Error,
		},
	}
	r := &mutationResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.UpdateFeature(ctx, input, 1, nil, make([]entity.AuditAction, 0), []string{"domain"}, nil, nil)
			if !tt.wantErr(t, err, fmt.Sprintf("FeatureRole")) {
				return
			}
			assert.Equalf(t, tt.want, got, "FeatureRole")
		})
	}
}

func Test_mutationResolver_DeleteFeature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockedService := mock.NewMockAuditService(ctrl)
	gomock.InOrder(
		mockedService.EXPECT().DeleteFeature(ctx, gomock.Any()).Return(nil),
		mockedService.EXPECT().DeleteFeature(ctx, gomock.Any()).Return(errors.New("Unknown error")),
	)
	tests := []struct {
		name    string
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "success",
			want:    "Feature deleted successfully",
			wantErr: assert.NoError,
		},
		{
			name:    "fail",
			want:    "Failed to delete feature",
			wantErr: assert.Error,
		},
	}
	r := &mutationResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.DeleteFeature(ctx, 1, nil, make([]entity.AuditAction, 0), []string{"domain"}, nil, nil)
			if !tt.wantErr(t, err, fmt.Sprintf("FeatureRole")) {
				return
			}
			assert.Equalf(t, tt.want, got, "FeatureRole")
		})
	}
}

func TestQueryResolver_SearchFeature_Success(t *testing.T) {
	ctx := context.Background()

	res := &presenter.SearchFeature{
		CurrentPage: 1,
		Pages:       1,
		PerPage:     1,
		TotalItems:  1,
		Items: []*entity.Feature{
			{
				ID:              1,
				Name:            nil,
				ParentFeatureID: new(int),
			},
		},
	}

	psm.Configure("console-audit-bff", factory.Default())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedService := mock.NewMockAuditService(ctrl)
	mockedService.EXPECT().
		SearchFeature(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(res, nil)

	r := &queryResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}
	// test
	response, err := r.SearchFeature(ctx, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, 1, response.CurrentPage)
	require.Equal(t, 1, response.Pages)
	require.Equal(t, 1, response.PerPage)
	require.Equal(t, 1, response.TotalItems)
	require.Equal(t, 1, response.Features[0].ID)
}

func TestQueryResolver_SearchFeature_Error(t *testing.T) {
	ctx := context.Background()

	res := &presenter.SearchFeature{
		CurrentPage: 1,
		Pages:       1,
		PerPage:     1,
		TotalItems:  1,
		Items: []*entity.Feature{
			{
				ID:              1,
				Name:            nil,
				ParentFeatureID: new(int),
			},
		},
	}

	psm.Configure("console-audit-bff", factory.Default())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedService := mock.NewMockAuditService(ctrl)
	mockedService.EXPECT().
		SearchFeature(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(res, errors.New("err"))

	r := &queryResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}
	// test
	response, err := r.SearchFeature(ctx, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, "err", err.Error())
}

func TestQueryResolver_FindFeatureByID_Success(t *testing.T) {
	ctx := context.Background()

	res := &entity.Feature{
		ID:              1,
		Name:            nil,
		ParentFeatureID: new(int),
	}

	psm.Configure("console-audit-bff", factory.Default())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedService := mock.NewMockAuditService(ctrl)
	mockedService.EXPECT().
		FindFeatureByID(ctx, gomock.Any()).
		Return(res, nil)

	r := &queryResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}
	// test
	response, err := r.FindFeatureByID(ctx, 1)
	require.NoError(t, err)
	require.NotNil(t, response)

}

func TestQueryResolver_FindFeatureByID_Error(t *testing.T) {
	ctx := context.Background()

	res := &entity.Feature{
		ID:              1,
		Name:            nil,
		ParentFeatureID: new(int),
	}

	psm.Configure("console-audit-bff", factory.Default())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedService := mock.NewMockAuditService(ctrl)
	mockedService.EXPECT().
		FindFeatureByID(ctx, gomock.Any()).
		Return(res, errors.New("err"))

	r := &queryResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}
	// test
	response, err := r.FindFeatureByID(ctx, 1)
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, "err", err.Error())

}

func Test_mutationResolver_AttachRoleToUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockedService := mock.NewMockAuditService(ctrl)
	gomock.InOrder(
		mockedService.EXPECT().AttachRoleToUser(ctx, gomock.Any(), gomock.Any()).Return(nil),
		mockedService.EXPECT().AttachRoleToUser(ctx, gomock.Any(), gomock.Any()).Return(errors.New("Unknown error")),
	)
	tests := []struct {
		name    string
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "success",
			want:    "Role attached successfully",
			wantErr: assert.NoError,
		},
		{
			name:    "fail",
			want:    "Failed to attach user to role",
			wantErr: assert.Error,
		},
	}
	r := &mutationResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.AttachRoleToUser(ctx, 1, "pismo@pismo.io")
			if !tt.wantErr(t, err, fmt.Sprintf("AttachRoleToUser")) {
				return
			}
			assert.Equalf(t, tt.want, got, "AttachRoleToUser")
		})
	}
}

func Test_mutationResolver_DetachRoleToUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockedService := mock.NewMockAuditService(ctrl)
	gomock.InOrder(
		mockedService.EXPECT().DetachRoleToUser(ctx, gomock.Any(), gomock.Any()).Return(nil),
		mockedService.EXPECT().DetachRoleToUser(ctx, gomock.Any(), gomock.Any()).Return(errors.New("Unknown error")),
	)
	tests := []struct {
		name    string
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "success",
			want:    "Role detached successfully",
			wantErr: assert.NoError,
		},
		{
			name:    "fail",
			want:    "Failed to detach user to role",
			wantErr: assert.Error,
		},
	}
	r := &mutationResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.DetachRoleToUser(ctx, 1, "pismo@pismo.io")
			if !tt.wantErr(t, err, fmt.Sprintf("DetachRoleToUser")) {
				return
			}
			assert.Equalf(t, tt.want, got, "DetachRoleToUser")
		})
	}
}

func Test_mutationResolver_AttachRoleToFeature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockedService := mock.NewMockAuditService(ctrl)
	gomock.InOrder(
		mockedService.EXPECT().AttachRoleToFeature(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
		mockedService.EXPECT().AttachRoleToFeature(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Unknown error")),
	)
	tests := []struct {
		name    string
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "success",
			want:    "Feature attached successfully",
			wantErr: assert.NoError,
		},
		{
			name:    "fail",
			want:    "Failed to attach feature to role",
			wantErr: assert.Error,
		},
	}
	r := &mutationResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.AttachRoleToFeature(ctx, 1, 1, entity.ActionInput{Action: "write"})
			if !tt.wantErr(t, err, fmt.Sprintf("AttachRoleToFeature")) {
				return
			}
			assert.Equalf(t, tt.want, got, "AttachRoleToFeature")
		})
	}
}

func Test_mutationResolver_DetachRoleToFeature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockedService := mock.NewMockAuditService(ctrl)
	gomock.InOrder(
		mockedService.EXPECT().DetachRoleToFeature(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
		mockedService.EXPECT().DetachRoleToFeature(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Unknown error")),
	)
	tests := []struct {
		name    string
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "success",
			want:    "Feature detached successfully",
			wantErr: assert.NoError,
		},
		{
			name:    "fail",
			want:    "Failed to detach feature to role",
			wantErr: assert.Error,
		},
	}
	r := &mutationResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.DetachRoleToFeature(ctx, 1, 1, entity.ActionInput{Action: "write"})
			if !tt.wantErr(t, err, fmt.Sprintf("DetachRoleToFeature")) {
				return
			}
			assert.Equalf(t, tt.want, got, "DetachRoleToFeature")
		})
	}
}

func Test_mutationResolver_CreateEndpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	input := entity.EndpointInput{}
	mockedService := mock.NewMockAuditService(ctrl)
	gomock.InOrder(
		mockedService.EXPECT().CreateEndpoint(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
		mockedService.EXPECT().CreateEndpoint(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Duplicated endpoint")),
		mockedService.EXPECT().CreateEndpoint(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Unknown error")),
	)
	tests := []struct {
		name    string
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "success",
			want:    "Endpoint created successfully",
			wantErr: assert.NoError,
		},
		{
			name:    "bad request",
			want:    "Duplicated endpoint",
			wantErr: assert.NoError,
		},
		{
			name:    "fail",
			want:    "Failed to create endpoint",
			wantErr: assert.Error,
		},
	}
	r := &mutationResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.CreateEndpoint(ctx, input, nil, make([]entity.AuditAction, 0), []string{"domain"}, nil, nil)
			if !tt.wantErr(t, err, fmt.Sprintf("CreateEndpoint(%v, %v)", ctx, input)) {
				return
			}
			assert.Equalf(t, tt.want, got, "CreateEndpoint(%v, %v)", ctx, input)
		})
	}
}

func Test_mutationResolver_UpdateEndpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	input := entity.EndpointInput{}
	mockedService := mock.NewMockAuditService(ctrl)
	gomock.InOrder(
		mockedService.EXPECT().UpdateEndpoint(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
		mockedService.EXPECT().UpdateEndpoint(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Duplicated endpoint")),
		mockedService.EXPECT().UpdateEndpoint(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Unknown error")),
	)
	tests := []struct {
		name    string
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "success",
			want:    "Endpoint updated successfully",
			wantErr: assert.NoError,
		},
		{
			name:    "bad request",
			want:    "Duplicated endpoint",
			wantErr: assert.NoError,
		},
		{
			name:    "fail",
			want:    "Failed to update endpoint",
			wantErr: assert.Error,
		},
	}
	r := &mutationResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.UpdateEndpoint(ctx, input, 1, nil, make([]entity.AuditAction, 0), []string{"domain"}, nil, nil)
			if !tt.wantErr(t, err, fmt.Sprintf("EndpointRole")) {
				return
			}
			assert.Equalf(t, tt.want, got, "EndpointRole")
		})
	}
}

func Test_mutationResolver_DeleteEndpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockedService := mock.NewMockAuditService(ctrl)
	gomock.InOrder(
		mockedService.EXPECT().DeleteEndpoint(ctx, gomock.Any()).Return(nil),
		mockedService.EXPECT().DeleteEndpoint(ctx, gomock.Any()).Return(errors.New("Unknown error")),
	)
	tests := []struct {
		name    string
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "success",
			want:    "Endpoint deleted successfully",
			wantErr: assert.NoError,
		},
		{
			name:    "fail",
			want:    "Failed to delete endpoint",
			wantErr: assert.Error,
		},
	}
	r := &mutationResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.DeleteEndpoint(ctx, 1, nil, make([]entity.AuditAction, 0), []string{"domain"}, nil, nil)
			if !tt.wantErr(t, err, fmt.Sprintf("EndpointRole")) {
				return
			}
			assert.Equalf(t, tt.want, got, "EndpointRole")
		})
	}
}

func TestQueryResolver_SearchEndpoint_Success(t *testing.T) {
	ctx := context.Background()

	res := &presenter.SearchEndpoint{
		CurrentPage: 1,
		Pages:       1,
		PerPage:     1,
		TotalItems:  1,
		Items: []*entity.Endpoint{
			{
				ID:          1,
				ServiceName: pointer("api-x"),
				Path:        pointer("/users/"),
				Method:      pointer("POST"),
			},
		},
	}

	psm.Configure("console-audit-bff", factory.Default())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedService := mock.NewMockAuditService(ctrl)
	mockedService.EXPECT().
		SearchEndpoint(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(res, nil)

	r := &queryResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}
	// test
	response, err := r.SearchEndpoint(ctx, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, 1, response.CurrentPage)
	require.Equal(t, 1, response.Pages)
	require.Equal(t, 1, response.PerPage)
	require.Equal(t, 1, response.TotalItems)
	require.Equal(t, 1, response.Endpoints[0].ID)
}

func TestQueryResolver_SearchEndpoint_Error(t *testing.T) {
	ctx := context.Background()

	res := &presenter.SearchEndpoint{
		CurrentPage: 1,
		Pages:       1,
		PerPage:     1,
		TotalItems:  1,
		Items: []*entity.Endpoint{
			{
				ID:          1,
				ServiceName: pointer("api-x"),
				Path:        pointer("/users/"),
				Method:      pointer("POST"),
			},
		},
	}

	psm.Configure("console-audit-bff", factory.Default())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedService := mock.NewMockAuditService(ctrl)
	mockedService.EXPECT().
		SearchEndpoint(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(res, errors.New("err"))

	r := &queryResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}
	// test
	response, err := r.SearchEndpoint(ctx, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, "err", err.Error())
}

func TestQueryResolver_FindEndpointByID_Success(t *testing.T) {
	ctx := context.Background()

	res := &entity.Endpoint{
		ID:          1,
		ServiceName: pointer("api-x"),
		Path:        pointer("/users/"),
		Method:      pointer("POST"),
	}

	psm.Configure("console-audit-bff", factory.Default())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedService := mock.NewMockAuditService(ctrl)
	mockedService.EXPECT().
		FindEndpointByID(ctx, gomock.Any()).
		Return(res, nil)

	r := &queryResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}
	// test
	response, err := r.FindEndpointByID(ctx, 1)
	require.NoError(t, err)
	require.NotNil(t, response)

}

func TestQueryResolver_FindEndpointByID_Error(t *testing.T) {
	ctx := context.Background()

	res := &entity.Endpoint{
		ID:          1,
		ServiceName: pointer("api-x"),
		Path:        pointer("/users/"),
		Method:      pointer("POST"),
	}

	psm.Configure("console-audit-bff", factory.Default())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedService := mock.NewMockAuditService(ctrl)
	mockedService.EXPECT().
		FindEndpointByID(ctx, gomock.Any()).
		Return(res, errors.New("err"))

	r := &queryResolver{
		Resolver: &Resolver{
			Manager: &service.Manager{
				AuditService: mockedService,
			},
		},
	}
	// test
	response, err := r.FindEndpointByID(ctx, 1)
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, "err", err.Error())

}

