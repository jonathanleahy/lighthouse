package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/pismo/backoffice-core-bff/internal/app/domain/service"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	res := httptest.NewRecorder()
	echoContext := echo.New().NewContext(req, res)

	SetResolver(&service.Resolver{HealthService: service.NewHealthService()})
	err := NewHealthHandler().Get(echoContext)

	assert.Nil(t, err)
}

