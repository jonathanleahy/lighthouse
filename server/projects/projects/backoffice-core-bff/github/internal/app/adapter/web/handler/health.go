package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HealthHandler interface {
	Get(c echo.Context) error
}

type healthHandler struct{}

func NewHealthHandler() HealthHandler {
	return &healthHandler{}
}

// Get godoc
// @Summary Return service status
// @Description Return service status
// @Tags health
// @Produce json
// @success 200 {object} presenter.Health
// @success 500 {object} presenter.Health
// @Router /health [get]
func (h healthHandler) Get(c echo.Context) error {
	health, _ := resolver.HealthService.GetMessage()
	return c.JSON(http.StatusOK, health)
}

