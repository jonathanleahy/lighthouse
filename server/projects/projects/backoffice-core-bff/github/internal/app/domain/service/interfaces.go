package service

import (
	"context"
	"github.com/pismo/backoffice-core-bff/internal/app/domain/entity"
)

type AuditService interface {
	SendAuditEvent(ctx context.Context, domain, action, domainID string, request string, response string, responseCode int) chan error
	WaitFinish()
}

type DisputeService interface {
	GetDisputeStatus(ctx context.Context, disputeId int, disputeInstallmentId *int) ([]*entity.DisputeStatus, error)
	CreateFraudReport(ctx context.Context, disputeId int, input entity.FraudReportInput) (*entity.FraudReport, error)
	UpdateDisputeStatusEvent(ctx context.Context, disputeId int, input entity.DisputeEventInput) (*entity.DisputeEvent, error)
}

type eventSender interface {
	SendEvent(ctx context.Context, domain string, eventType string, body interface{}) error
}

type HealthService interface {
	GetMessage() (*entity.Health, error)
}

