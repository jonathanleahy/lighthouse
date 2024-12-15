package utils

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/request"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetRequestIP(t *testing.T) {
	type args struct {
		xFowardedFor string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should return a valid request IP",
			args: args{"191.168.1.1,192.168.1.2,192.168.1.3"},
			want: "191.168.1.1",
		},
		{
			name: "should return a blank request IP",
			args: args{""},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRequestIP(tt.args.xFowardedFor); got != tt.want {
				t.Errorf("GetRequestIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCIDFromContext_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{})

	result := GetCIDFromContext(ctx)

	assert.NotNil(t, result)
}

func TestGetTenantFromContext_Success(t *testing.T) {
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{})

	result := GetTenantFromContext(ctx)

	assert.NotNil(t, result)
}

func TestGetCIDFromEchoContext_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	res := httptest.NewRecorder()
	echoContext := echo.New().NewContext(req, res)

	result := GetCIDFromEchoContext(echoContext)

	assert.NotNil(t, result)
}

func TestGetTenantFromEchoContext_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	res := httptest.NewRecorder()
	echoContext := echo.New().NewContext(req, res)

	result := GetTenantFromEchoContext(echoContext)

	assert.NotNil(t, result)
}

func TestGetRequestBody_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	res := httptest.NewRecorder()
	echoContext := echo.New().NewContext(req, res)

	result := GetRequestBody(echoContext)

	assert.NotNil(t, result)
}

