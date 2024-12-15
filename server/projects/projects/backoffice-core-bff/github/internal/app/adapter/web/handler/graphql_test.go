package handler

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/pismo/backoffice-core-bff/internal/app/domain/service"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/apierror"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/request"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestGetGraphqlServer(t *testing.T) {
	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

}

func Test_verifyQueryIsAllowed(t *testing.T) {
	tests := []struct {
		name string
		args echo.Context
		want error
	}{
		{
			name: "should validate query without schema",
			args: buildQueryWhiteListContext(jsonRequestWithoutSchema),
			want: nil,
		},
		{
			name: "should validate query without query field",
			args: buildQueryWhiteListContext(jsonRequestWithoutQueryField),
			want: nil,
		},
	}
	_ = os.Setenv(env.Env, "test")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := verifyQueryIsAllowed(tt.args); got != tt.want {
				t.Errorf("verifyQueryIsAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyQueryIsAllowed_Json_Error(t *testing.T) {
	_ = os.Setenv(env.Env, "test")
	ctx := buildQueryWhiteListContext(jsonRequestWithWrongJson)
	err := verifyQueryIsAllowed(ctx)

	assert.NotNil(t, err)
}

func TestVerifyQueryIsAllowedWithDevEnv(t *testing.T) {
	_ = os.Setenv(env.Env, "dev")
	err := verifyQueryIsAllowed(nil)

	assert.Nil(t, err)
}

func TestVerifyQueryIsAllowedWithSchema(t *testing.T) {
	_ = os.Setenv(env.Env, "prod")
	ctx := buildQueryWhiteListContext(jsonRequestWithSchema)
	err := verifyQueryIsAllowed(ctx)

	apiError := err.(apierror.ApiError)
	assert.Equal(t, "ECPBFF0014", apiError.GetMessage().ErrorCode)
	assert.Equal(t, "Query not allowed", apiError.GetMessage().UserMessage)
}

func buildQueryWhiteListContext(json string) echo.Context {
	rctx := request.RequestContext{
		Tenant:        "X-TENANT-X",
		Cid:           "123",
		CustomHeaders: make(map[string]string),
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := context.WithValue(c.Request().Context(), env.RequestContext, rctx)
	c.SetRequest(c.Request().WithContext(ctx))
	return c
}

func TestGetGraphqlQuery(t *testing.T) {
	_ = os.Setenv(env.Env, "prod")
	SetResolver(&service.Resolver{GraphqlService: service.NewGraphqlService(&service.Resolver{})})

	tests := []struct {
		name string
		arg  string
	}{
		{name: "Status Forbidden", arg: jsonRequestWithSchema},
		{name: "Serve HTTP", arg: jsonRequestWithoutSchema},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewGraphqlHandler()
			ctx := buildQueryWhiteListContext(tt.arg)
			err := handler.Query(ctx)

			assert.NotNil(t, handler)
			assert.Nil(t, err)
		})
	}
}

func TestGetGraphqlPlayground(t *testing.T) {
	_ = os.Setenv(env.Env, "prod")
	SetResolver(&service.Resolver{GraphqlService: service.NewGraphqlService(&service.Resolver{})})

	handler := NewGraphqlHandler()
	ctx := buildQueryWhiteListContext(jsonRequestWithSchema)
	err := handler.Playground(ctx)

	assert.NotNil(t, handler)
	assert.Nil(t, err)
}

var (
	jsonRequestWithSchema        = `{"operationName":null,"variables":{},"query":"query MyQuery { \n  __schema { \n    types { \n      name \n      fields { \n      name \n      } \n    } \n  } \n}"}`
	jsonRequestWithoutSchema     = `{"operationName":null,"variables":{},"query":"query MyQuery { \n  program { \n    types { \n      name \n      fields { \n      name \n      } \n    } \n  } \n}"}`
	jsonRequestWithoutQueryField = `{"operationName":null,"variables":{},"other":"query MyQuery { \n  program { \n    types { \n      name \n      fields { \n      name \n      } \n    } \n  } \n}"}`
	jsonRequestWithWrongJson     = `{{}`
)

