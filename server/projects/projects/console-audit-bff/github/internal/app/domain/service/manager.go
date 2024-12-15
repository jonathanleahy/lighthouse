package service

import (
	"github.com/pismo/psm-sdk/psm/env"
	"github.com/pismo/psm-sdk/psm/schema"
)

type Manager struct {
	AuditService AuditService
	NetService   NetService
}

func NewManager(httpClient schema.HttpClient) *Manager {
	return &Manager{
		AuditService: NewAuditService(httpClient, env.Get("CONSOLE_AUDIT_API_URL", "https://console-audit-api.integration.pismolabs.io:443")),
		NetService:   NewNetService(),
	}
}

