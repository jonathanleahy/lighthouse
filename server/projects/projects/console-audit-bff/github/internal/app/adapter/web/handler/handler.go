package handler

import (
	"context"
	"net/http"

	"github.com/pismo/console-audit-bff/internal/app/domain/service"
	"github.com/pismo/psm-sdk/psm/server"
)

var manager *service.Manager

func SetManager(m *service.Manager) {
	manager = m
}

func Whoami(ctx context.Context, req server.Request, res server.Response) error {

	ip := manager.NetService.Whoami(ctx, req.EchoCtx())

	return res.JSON(http.StatusOK, ip)
}

