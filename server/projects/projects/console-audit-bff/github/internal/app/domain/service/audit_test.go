package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/pismo/psm-sdk/psm/schema"

	"github.com/golang/mock/gomock"
	"github.com/pismo/console-audit-bff/internal/app/adapter/web/graphql/entity"
	"github.com/pismo/console-audit-bff/internal/app/adapter/web/presenter"
	"github.com/pismo/psm-sdk/psm"
	"github.com/pismo/psm-sdk/psm/attribute/value"
	"github.com/pismo/psm-sdk/psm/factory"
	"github.com/pismo/psm-sdk/psm/network/http"
	mockhttp "github.com/pismo/psm-sdk/psm/schema/mock"
	"github.com/pismo/psm-sdk/psm/server/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuditService(t *testing.T) {
	psm.Configure("console-audit-bff", factory.Default())

	assert.Equal(t, &auditService{
		httpClient: psm.HttpClient(),
		apiUrl:     "",
	}, NewAuditService(psm.HttpClient(), ""))
}

func TestAuditService_FindAuditByID(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		res        *presenter.Audit
		err        error
	}{
		{
			name:       "Success",
			statusCode: http.StatusOK,
			res: &presenter.Audit{
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
				Http: &presenter.Http{
					Code:     200,
					Request:  `{"name": "crm-operator", "email": "crm-operator@pismo.io", "roles": ["crm-operator"]}`,
					Response: `{"id": 1}`,
				},
			},
		},
		{
			name:       "ErrorNotFound",
			statusCode: http.StatusNotFound,
			err:        errors.New("not_found"),
		},
		{
			name:       "Error",
			statusCode: http.StatusInternalServerError,
			err:        errors.New("connection_err"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockedHttpClient := mockhttp.NewMockHttpClient(ctrl)
			apiUrl := "http://localhost:8080"

			srv := NewAuditService(mockedHttpClient, apiUrl)
			require.NotNil(t, srv)

			// creates a context with request data
			ctx := request.Inject(context.Background(), &request.Context{
				Tenant: "TN-xxx",
				Email:  "owner@pismo.io",
			})

			// creates the expected http request
			httpReq := &http.Request{
				Method: "GET",
				Host:   apiUrl,
				Route:  http.Route("/v1/audit/:id").PathParam("id", value.Plain(1)),
				Span:   http.Span().Attr("id", value.Plain(1)),
			}

			// creates a mocked http response
			httpRes := http.NewResponse(testCase.statusCode, nil, nil, testCase.res, testCase.err)

			// assert order execution using gomock library
			gomock.InOrder(
				mockedHttpClient.EXPECT().Execute(ctx, httpReq, gomock.Any()).Return(httpRes),
			)

			// executes the operation to be tested
			response, err := srv.FindAuditByID(ctx, 1)
			if testCase.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, response, testCase.res)
			} else {
				assert.Equal(t, testCase.err, testCase.err)
				assert.Equal(t, response, testCase.res)
			}
		})
	}
}

func TestAuditService_SearchAudit(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		res        *presenter.SearchAudit
		err        error
	}{
		{
			name:       "Success",
			statusCode: http.StatusOK,
			res: &presenter.SearchAudit{
				CurrentPage: 1,
				PerPage:     1,
				Pages:       1,
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
						Http: &presenter.Http{
							Code:     200,
							Request:  `{"name": "crm-operator", "email": "crm-operator@pismo.io", "roles": ["crm-operator"]}`,
							Response: `{"id": 1}`,
						},
					},
				},
			},
		},
		{
			name:       "Error",
			statusCode: http.StatusInternalServerError,
			err:        errors.New("connection_err"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockedHttpClient := mockhttp.NewMockHttpClient(ctrl)
			apiUrl := "http://localhost:8080"

			srv := NewAuditService(mockedHttpClient, apiUrl)
			require.NotNil(t, srv)

			ctx := context.TODO()

			// creates a mocked http response
			httpRes := http.NewResponse(testCase.statusCode, nil, nil, testCase.res, testCase.err)

			// assert order execution using gomock library
			gomock.InOrder(
				mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(httpRes),
			)

			order := new(entity.Order)
			*order = entity.OrderAsc

			var actions []entity.AuditAction
			actions = append(actions, entity.AuditActionCreate)

			// executes the operation to be tested
			response, err := srv.SearchAudit(ctx, pointer(1), pointer(1), order, pointer("2020-10-10"), pointer("2020-10-10"), pointer("email"), actions, []string{"domain"}, pointer("domain_id"), pointer("cid"))
			if testCase.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, response, testCase.res)
			} else {
				assert.Equal(t, testCase.err, testCase.err)
				assert.Equal(t, response, testCase.res)
			}
		})
	}
}

func TestRoleService_ListRoles(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		res        []*entity.Role
		err        error
	}{
		{
			name:       "Success",
			statusCode: http.StatusOK,
			res:        []*entity.Role{},
		},
		{
			name:       "Error",
			statusCode: http.StatusInternalServerError,
			err:        errors.New("connection_err"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockedHttpClient := mockhttp.NewMockHttpClient(ctrl)
			apiUrl := "http://localhost:8080"

			srv := NewAuditService(mockedHttpClient, apiUrl)
			require.NotNil(t, srv)

			ctx := context.TODO()

			// creates a mocked http response
			httpRes := http.NewResponse(testCase.statusCode, nil, nil, testCase.res, testCase.err)

			// assert order execution using gomock library
			gomock.InOrder(
				mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(httpRes),
			)

			var actions []entity.AuditAction
			actions = append(actions, entity.AuditActionCreate)

			// executes the operation to be tested
			response, err := srv.ListRoles(ctx)
			if testCase.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, response, testCase.res)
			} else {
				assert.Equal(t, testCase.err, testCase.err)
				assert.Equal(t, response, testCase.res)
			}
		})
	}
}

func TestRoleService_ListUserRoles(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		res        *presenter.UserRoles
		err        error
	}{
		{
			name:       "Success",
			statusCode: http.StatusOK,
			res: &presenter.UserRoles{
				Roles: []string{"admin", "crm-operator"},
			},
		},
		{
			name:       "Error",
			statusCode: http.StatusInternalServerError,
			err:        errors.New("connection_err"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockedHttpClient := mockhttp.NewMockHttpClient(ctrl)
			apiUrl := "http://localhost:8080"

			srv := NewAuditService(mockedHttpClient, apiUrl)
			require.NotNil(t, srv)

			ctx := context.TODO()

			// creates a mocked http response
			httpRes := http.NewResponse(testCase.statusCode, nil, nil, testCase.res, testCase.err)

			// assert order execution using gomock library
			gomock.InOrder(
				mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(httpRes),
			)

			var actions []entity.AuditAction
			actions = append(actions, entity.AuditActionCreate)

			// executes the operation to be tested
			response, err := srv.ListUserRoles(ctx, pointer("email"), actions, []string{"domain"}, pointer("domain_id"), pointer("cid"))
			if testCase.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, response, testCase.res)
			} else {
				assert.Equal(t, testCase.err, testCase.err)
				assert.Equal(t, response, testCase.res)
			}
		})
	}
}

func pointer[T any](value T) *T {
	return &value
}

func TestRoleService_SearchRoles(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		res        presenter.Roles
		err        error
	}{
		{
			name:       "Success",
			statusCode: http.StatusOK,
			res: presenter.Roles{
				"organization": map[string]map[string][]string{"access_key": {"read": []string{"admin"}, "view": []string{"admin", "crm-operator"}}},
			},
		},
		{
			name:       "Error",
			statusCode: http.StatusInternalServerError,
			err:        errors.New("connection_err"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ctx := context.Background()
			mockedHttpClient := mockhttp.NewMockHttpClient(ctrl)
			srv := NewAuditService(mockedHttpClient, "http://localhost:8080")
			// creates a mocked http response
			httpRes := http.NewResponse(testCase.statusCode, nil, nil, testCase.res, testCase.err)
			mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(httpRes)

			// executes the operation to be tested
			response, err := srv.SearchRoles(ctx, pointer("email"), pointer("feature"))
			if testCase.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, response, testCase.res)
			} else {
				assert.Equal(t, testCase.err, testCase.err)
				assert.Equal(t, response, testCase.res)
			}
		})
	}
}

func Test_auditService_CreateRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	input := &entity.RoleInput{}
	mockedHttpClient := mockhttp.NewMockHttpClient(ctrl)
	gomock.InOrder(
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusOK, nil, nil, nil, nil),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusBadRequest, nil, nil, nil, errors.New("error")),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusInternalServerError, nil, nil, nil, errors.New("error")),
		),
	)

	tests := []struct {
		name       string
		httpClient schema.HttpClient
		wantErr    assert.ErrorAssertionFunc
	}{
		{name: "Success", wantErr: assert.NoError},
		{name: "Bad Request", wantErr: assert.Error},
		{name: "Fail", wantErr: assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &auditService{
				httpClient: mockedHttpClient,
				apiUrl:     "",
			}
			tt.wantErr(t, s.CreateRole(ctx, input), fmt.Sprintf("CreateRole(%v, %v)", ctx, input))
		})
	}
}

func Test_auditService_UpdateRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	input := &entity.RoleInput{}
	mockedHttpClient := mockhttp.NewMockHttpClient(ctrl)
	gomock.InOrder(
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusOK, nil, nil, nil, nil),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusBadRequest, nil, nil, nil, errors.New("error")),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusInternalServerError, nil, nil, nil, errors.New("error")),
		),
	)

	tests := []struct {
		name       string
		httpClient schema.HttpClient
		wantErr    assert.ErrorAssertionFunc
	}{
		{name: "Success", wantErr: assert.NoError},
		{name: "Bad Request", wantErr: assert.Error},
		{name: "Fail", wantErr: assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &auditService{
				httpClient: mockedHttpClient,
				apiUrl:     "",
			}
			tt.wantErr(t, s.UpdateRole(ctx, input), fmt.Sprintf("CreateRole(%v, %v)", ctx, input))
		})
	}
}

func Test_auditService_DeleteRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockedHttpClient := mockhttp.NewMockHttpClient(ctrl)
	gomock.InOrder(
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusOK, nil, nil, nil, nil),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusBadRequest, nil, nil, nil, errors.New("error")),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusInternalServerError, nil, nil, nil, errors.New("error")),
		),
	)

	tests := []struct {
		name       string
		httpClient schema.HttpClient
		wantErr    assert.ErrorAssertionFunc
	}{
		{name: "Success", wantErr: assert.NoError},
		{name: "Bad Request", wantErr: assert.Error},
		{name: "Fail", wantErr: assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &auditService{
				httpClient: mockedHttpClient,
				apiUrl:     "",
			}
			tt.wantErr(t, s.DeleteRole(ctx, 1), fmt.Sprintf("CreateRole(%v, %v)", ctx, 1))
		})
	}
}

func Test_auditService_CreateFeature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	input := &entity.FeatureInput{}
	mockedHttpClient := mockhttp.NewMockHttpClient(ctrl)
	gomock.InOrder(
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusOK, nil, nil, nil, nil),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusBadRequest, nil, nil, nil, errors.New("error")),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusInternalServerError, nil, nil, nil, errors.New("error")),
		),
	)

	tests := []struct {
		name       string
		httpClient schema.HttpClient
		wantErr    assert.ErrorAssertionFunc
	}{
		{name: "Success", wantErr: assert.NoError},
		{name: "Bad Request", wantErr: assert.Error},
		{name: "Fail", wantErr: assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &auditService{
				httpClient: mockedHttpClient,
				apiUrl:     "",
			}
			tt.wantErr(t, s.CreateFeature(ctx, input, nil, make([]entity.AuditAction, 0), []string{"foobar"}, nil, nil), fmt.Sprintf("CreateFeature(%v, %v)", ctx, input))
		})
	}
}

func Test_auditService_UpdateFeature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	input := &entity.FeatureInput{}
	mockedHttpClient := mockhttp.NewMockHttpClient(ctrl)
	gomock.InOrder(
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusOK, nil, nil, nil, nil),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusBadRequest, nil, nil, nil, errors.New("error")),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusInternalServerError, nil, nil, nil, errors.New("error")),
		),
	)

	tests := []struct {
		name       string
		httpClient schema.HttpClient
		wantErr    assert.ErrorAssertionFunc
	}{
		{name: "Success", wantErr: assert.NoError},
		{name: "Bad Request", wantErr: assert.Error},
		{name: "Fail", wantErr: assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &auditService{
				httpClient: mockedHttpClient,
				apiUrl:     "",
			}
			tt.wantErr(t, s.UpdateFeature(ctx, input, 1, nil, make([]entity.AuditAction, 0), []string{"foobar"}, nil, nil), fmt.Sprintf("UpdateFeature(%v, %v)", ctx, input))
		})
	}
}

func Test_auditService_DeleteFeature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	input := &entity.FeatureInput{}
	mockedHttpClient := mockhttp.NewMockHttpClient(ctrl)
	gomock.InOrder(
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusOK, nil, nil, nil, nil),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusBadRequest, nil, nil, nil, errors.New("error")),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusInternalServerError, nil, nil, nil, errors.New("error")),
		),
	)

	tests := []struct {
		name       string
		httpClient schema.HttpClient
		wantErr    assert.ErrorAssertionFunc
	}{
		{name: "Success", wantErr: assert.NoError},
		{name: "Bad Request", wantErr: assert.Error},
		{name: "Fail", wantErr: assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &auditService{
				httpClient: mockedHttpClient,
				apiUrl:     "",
			}
			tt.wantErr(t, s.DeleteFeature(ctx, 1), fmt.Sprintf("DeleteFeature(%v, %v)", ctx, input))
		})
	}
}

func TestAuditService_FindFeatureByID(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		res        *entity.Feature
		err        error
	}{
		{
			name:       "Success",
			statusCode: http.StatusOK,
			res: &entity.Feature{
				ID: 1,
			},
		},
		{
			name:       "ErrorNotFound",
			statusCode: http.StatusNotFound,
			err:        errors.New("not_found"),
		},
		{
			name:       "Error",
			statusCode: http.StatusInternalServerError,
			err:        errors.New("connection_err"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockedHttpClient := mockhttp.NewMockHttpClient(ctrl)
			apiUrl := "http://localhost:8080"

			srv := NewAuditService(mockedHttpClient, apiUrl)
			require.NotNil(t, srv)

			// creates a context with request data
			ctx := request.Inject(context.Background(), &request.Context{
				Tenant: "TN-xxx",
				Email:  "owner@pismo.io",
			})

			// creates the expected http request
			httpReq := &http.Request{
				Method: "GET",
				Host:   apiUrl,
				Route:  http.Route("/v1/feature/:id").PathParam("id", value.Plain(1)),
				Span:   http.Span().Attr("id", value.Plain(1)),
			}

			// creates a mocked http response
			httpRes := http.NewResponse(testCase.statusCode, nil, nil, testCase.res, testCase.err)

			// assert order execution using gomock library
			gomock.InOrder(
				mockedHttpClient.EXPECT().Execute(ctx, httpReq, gomock.Any()).Return(httpRes),
			)

			// executes the operation to be tested
			response, err := srv.FindFeatureByID(ctx, 1)
			if testCase.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, response, testCase.res)
			} else {
				assert.Equal(t, testCase.err, testCase.err)
				assert.Equal(t, response, testCase.res)
			}
		})
	}
}

func TestAuditService_SearchFeature(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		res        *presenter.SearchFeature
		err        error
	}{
		{
			name:       "Success",
			statusCode: http.StatusOK,
			res: &presenter.SearchFeature{
				CurrentPage: 1,
				PerPage:     1,
				Pages:       1,
				TotalItems:  1,
				Items: []*entity.Feature{
					{
						ID: 1,
					},
				},
			},
		},
		{
			name:       "Error",
			statusCode: http.StatusInternalServerError,
			err:        errors.New("connection_err"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockedHttpClient := mockhttp.NewMockHttpClient(ctrl)
			apiUrl := "http://localhost:8080"

			srv := NewAuditService(mockedHttpClient, apiUrl)
			require.NotNil(t, srv)

			ctx := context.TODO()

			// creates a mocked http response
			httpRes := http.NewResponse(testCase.statusCode, nil, nil, testCase.res, testCase.err)

			// assert order execution using gomock library
			gomock.InOrder(
				mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(httpRes),
			)

			order := new(entity.Order)
			*order = entity.OrderAsc

			var actions []entity.AuditAction
			actions = append(actions, entity.AuditActionCreate)

			response, err := srv.SearchFeature(ctx, pointer(1), pointer(1), order, pointer("id"), pointer("name"), pointer("1"), pointer("user"), actions, []string{"domains"}, pointer("domain-id"), pointer("cid"))
			if testCase.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, response, testCase.res)
			} else {
				assert.Equal(t, testCase.err, testCase.err)
				assert.Equal(t, response, testCase.res)
			}
		})
	}
}
func Test_auditService_AttachRoleToUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockedHttpClient := mockhttp.NewMockHttpClient(ctrl)
	gomock.InOrder(
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusOK, nil, nil, nil, nil),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusBadRequest, nil, nil, nil, errors.New("error")),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusInternalServerError, nil, nil, nil, errors.New("error")),
		),
	)

	tests := []struct {
		name       string
		httpClient schema.HttpClient
		wantErr    assert.ErrorAssertionFunc
	}{
		{name: "Success", wantErr: assert.NoError},
		{name: "Bad Request", wantErr: assert.Error},
		{name: "Fail", wantErr: assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &auditService{
				httpClient: mockedHttpClient,
				apiUrl:     "",
			}
			tt.wantErr(t, s.AttachRoleToUser(ctx, 1, "pismo@pismo.io"), fmt.Sprintf("AttachRoleToUser(%v, %v)", ctx, 1))
		})
	}
}

func Test_auditService_DetachRoleToUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockedHttpClient := mockhttp.NewMockHttpClient(ctrl)
	gomock.InOrder(
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusOK, nil, nil, nil, nil),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusBadRequest, nil, nil, nil, errors.New("error")),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusInternalServerError, nil, nil, nil, errors.New("error")),
		),
	)

	tests := []struct {
		name       string
		httpClient schema.HttpClient
		wantErr    assert.ErrorAssertionFunc
	}{
		{name: "Success", wantErr: assert.NoError},
		{name: "Bad Request", wantErr: assert.Error},
		{name: "Fail", wantErr: assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &auditService{
				httpClient: mockedHttpClient,
				apiUrl:     "",
			}
			tt.wantErr(t, s.DetachRoleToUser(ctx, 1, "pismo@pismo.io"), fmt.Sprintf("DetachRoleToUser(%v, %v)", ctx, 1))
		})
	}
}

func Test_auditService_AttachRoleToFeature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockedHttpClient := mockhttp.NewMockHttpClient(ctrl)
	gomock.InOrder(
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusOK, nil, nil, nil, nil),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusBadRequest, nil, nil, nil, errors.New("error")),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusInternalServerError, nil, nil, nil, errors.New("error")),
		),
	)

	tests := []struct {
		name       string
		httpClient schema.HttpClient
		wantErr    assert.ErrorAssertionFunc
	}{
		{name: "Success", wantErr: assert.NoError},
		{name: "Bad Request", wantErr: assert.Error},
		{name: "Fail", wantErr: assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &auditService{
				httpClient: mockedHttpClient,
				apiUrl:     "",
			}
			tt.wantErr(t, s.AttachRoleToFeature(ctx, 1, 1, entity.ActionInput{Action: "write"}), fmt.Sprintf("AttachRoleToFeature(%v, %v)", ctx, 1))
		})
	}
}

func Test_auditService_DetachRoleToFeature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockedHttpClient := mockhttp.NewMockHttpClient(ctrl)
	gomock.InOrder(
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusOK, nil, nil, nil, nil),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusBadRequest, nil, nil, nil, errors.New("error")),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusInternalServerError, nil, nil, nil, errors.New("error")),
		),
	)

	tests := []struct {
		name       string
		httpClient schema.HttpClient
		wantErr    assert.ErrorAssertionFunc
	}{
		{name: "Success", wantErr: assert.NoError},
		{name: "Bad Request", wantErr: assert.Error},
		{name: "Fail", wantErr: assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &auditService{
				httpClient: mockedHttpClient,
				apiUrl:     "",
			}
			tt.wantErr(t, s.DetachRoleToFeature(ctx, 1, 1, entity.ActionInput{Action: "write"}), fmt.Sprintf("DetachRoleToFeature(%v, %v)", ctx, 1))
		})
	}
}

func Test_auditService_CreateEndpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	input := &entity.EndpointInput{}
	mockedHttpClient := mockhttp.NewMockHttpClient(ctrl)
	gomock.InOrder(
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusOK, nil, nil, nil, nil),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusBadRequest, nil, nil, nil, errors.New("error")),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusInternalServerError, nil, nil, nil, errors.New("error")),
		),
	)

	tests := []struct {
		name       string
		httpClient schema.HttpClient
		wantErr    assert.ErrorAssertionFunc
	}{
		{name: "Success", wantErr: assert.NoError},
		{name: "Bad Request", wantErr: assert.Error},
		{name: "Fail", wantErr: assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &auditService{
				httpClient: mockedHttpClient,
				apiUrl:     "",
			}
			tt.wantErr(t, s.CreateEndpoint(ctx, input, nil, make([]entity.AuditAction, 0), []string{"foobar"}, nil, nil), fmt.Sprintf("CreateEndpoint(%v, %v)", ctx, input))
		})
	}
}

func Test_auditService_UpdateEndpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	input := &entity.EndpointInput{}
	mockedHttpClient := mockhttp.NewMockHttpClient(ctrl)
	gomock.InOrder(
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusOK, nil, nil, nil, nil),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusBadRequest, nil, nil, nil, errors.New("error")),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusInternalServerError, nil, nil, nil, errors.New("error")),
		),
	)

	tests := []struct {
		name       string
		httpClient schema.HttpClient
		wantErr    assert.ErrorAssertionFunc
	}{
		{name: "Success", wantErr: assert.NoError},
		{name: "Bad Request", wantErr: assert.Error},
		{name: "Fail", wantErr: assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &auditService{
				httpClient: mockedHttpClient,
				apiUrl:     "",
			}
			tt.wantErr(t, s.UpdateEndpoint(ctx, input, 1, nil, make([]entity.AuditAction, 0), []string{"foobar"}, nil, nil), fmt.Sprintf("UpdateEndpoint(%v, %v)", ctx, input))
		})
	}
}

func Test_auditService_DeleteEndpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	input := &entity.EndpointInput{}
	mockedHttpClient := mockhttp.NewMockHttpClient(ctrl)
	gomock.InOrder(
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusOK, nil, nil, nil, nil),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusBadRequest, nil, nil, nil, errors.New("error")),
		),
		mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(
			http.NewResponse(http.StatusInternalServerError, nil, nil, nil, errors.New("error")),
		),
	)

	tests := []struct {
		name       string
		httpClient schema.HttpClient
		wantErr    assert.ErrorAssertionFunc
	}{
		{name: "Success", wantErr: assert.NoError},
		{name: "Bad Request", wantErr: assert.Error},
		{name: "Fail", wantErr: assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &auditService{
				httpClient: mockedHttpClient,
				apiUrl:     "",
			}
			tt.wantErr(t, s.DeleteEndpoint(ctx, 1), fmt.Sprintf("DeleteEndpoint(%v, %v)", ctx, input))
		})
	}
}

func TestAuditService_FindEndpointByID(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		res        *entity.Endpoint
		err        error
	}{
		{
			name:       "Success",
			statusCode: http.StatusOK,
			res: &entity.Endpoint{
				ID: 1,
			},
		},
		{
			name:       "ErrorNotFound",
			statusCode: http.StatusNotFound,
			err:        errors.New("not_found"),
		},
		{
			name:       "Error",
			statusCode: http.StatusInternalServerError,
			err:        errors.New("connection_err"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockedHttpClient := mockhttp.NewMockHttpClient(ctrl)
			apiUrl := "http://localhost:8080"

			srv := NewAuditService(mockedHttpClient, apiUrl)
			require.NotNil(t, srv)

			// creates a context with request data
			ctx := request.Inject(context.Background(), &request.Context{
				Tenant: "TN-xxx",
				Email:  "owner@pismo.io",
			})

			// creates the expected http request
			httpReq := &http.Request{
				Method: "GET",
				Host:   apiUrl,
				Route:  http.Route("/v1/endpoint/:id").PathParam("id", value.Plain(1)),
				Span:   http.Span().Attr("id", value.Plain(1)),
			}

			// creates a mocked http response
			httpRes := http.NewResponse(testCase.statusCode, nil, nil, testCase.res, testCase.err)

			// assert order execution using gomock library
			gomock.InOrder(
				mockedHttpClient.EXPECT().Execute(ctx, httpReq, gomock.Any()).Return(httpRes),
			)

			// executes the operation to be tested
			response, err := srv.FindEndpointByID(ctx, 1)
			if testCase.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, response, testCase.res)
			} else {
				assert.Equal(t, testCase.err, testCase.err)
				assert.Equal(t, response, testCase.res)
			}
		})
	}
}

func TestAuditService_SearchEndpoint(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		res        *presenter.SearchEndpoint
		err        error
	}{
		{
			name:       "Success",
			statusCode: http.StatusOK,
			res: &presenter.SearchEndpoint{
				CurrentPage: 1,
				PerPage:     1,
				Pages:       1,
				TotalItems:  1,
				Items: []*entity.Endpoint{
					{
						ID: 1,
					},
				},
			},
		},
		{
			name:       "Error",
			statusCode: http.StatusInternalServerError,
			err:        errors.New("connection_err"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockedHttpClient := mockhttp.NewMockHttpClient(ctrl)
			apiUrl := "http://localhost:8080"

			srv := NewAuditService(mockedHttpClient, apiUrl)
			require.NotNil(t, srv)

			ctx := context.TODO()

			// creates a mocked http response
			httpRes := http.NewResponse(testCase.statusCode, nil, nil, testCase.res, testCase.err)

			// assert order execution using gomock library
			gomock.InOrder(
				mockedHttpClient.EXPECT().Execute(ctx, gomock.Any(), gomock.Any()).Return(httpRes),
			)

			order := new(entity.Order)
			*order = entity.OrderAsc

			var actions []entity.AuditAction
			actions = append(actions, entity.AuditActionCreate)

			response, err := srv.SearchEndpoint(ctx, pointer(1), pointer(1), order, pointer("id"), pointer("service_name"), pointer("method"), pointer("path"), pointer("user"), actions, []string{"domains"}, pointer("domain-id"), pointer("cid"))
			if testCase.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, response, testCase.res)
			} else {
				assert.Equal(t, testCase.err, testCase.err)
				assert.Equal(t, response, testCase.res)
			}
		})
	}
}

