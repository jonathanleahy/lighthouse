package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
	//"github.com/pismo/psm-sdk/psm"
	//"github.com/pismo/psm-sdk/psm/env"
	//"github.com/pismo/psm-sdk/psm/factory"
	//"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	psm.Configure("console-audit-bff", factory.Default())
	expected := &Manager{
		AuditService: NewAuditService(psm.HttpClient(), env.Get("CONSOLE_AUDIT_API_URL", "https://console-audit-api.integration.pismolabs.io:443")),
		NetService:   NewNetService(),
	}

	assert.Equal(t, expected, NewManager(psm.HttpClient()))
}
