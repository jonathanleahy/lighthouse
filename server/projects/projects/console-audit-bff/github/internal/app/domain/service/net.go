package service

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/pismo/console-audit-bff/internal/app/adapter/web/presenter"
)

type (
	NetService interface {
		Whoami(ctx context.Context, netCtx echo.Context) *presenter.IP
	}

	netService struct {
	}
)

func NewNetService() NetService {
	return &netService{}
}

func (s *netService) Whoami(_ context.Context, netCtx echo.Context) *presenter.IP {
	return &presenter.IP{IP: netCtx.RealIP()}
}

