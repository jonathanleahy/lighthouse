package service

import (
	"context"
	"github.com/pismo/backoffice-core-bff/internal/app/domain/entity"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/request"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/utils"
	"strconv"
	"sync"
	"time"
)

type auditService struct {
	wg                    sync.WaitGroup
	snsConsoleAuditSender eventSender
}

func NewAuditService(snsConsoleAuditSender eventSender) AuditService {
	return &auditService{
		snsConsoleAuditSender: snsConsoleAuditSender,
	}
}

func (s *auditService) SendAuditEvent(ctx context.Context, domain, action, domainID string, request string, response string, responseCode int) chan error {
	s.wg.Add(1)

	chanErr := make(chan error, 1)
	go func() {
		defer s.wg.Done()
		chanErr <- s.sendAuditEvent(ctx, domain, action, domainID, request, response, responseCode)
	}()

	return chanErr
}

func (s *auditService) sendAuditEvent(ctx context.Context, domain, action, domainID string, request string, response string, responseCode int) error {
	s.wg.Add(1)
	defer s.wg.Done()

	audit := populateAudit(ctx, domain, action, domainID, request, response, responseCode)
	return s.snsConsoleAuditSender.SendEvent(ctx, domain, action, audit)
}

func (s *auditService) WaitFinish() {
	s.wg.Wait()
}

func populateAudit(ctx context.Context, domain, action, domainId, req string, response string, code int) entity.Audit {
	rctx := ctx.Value(env.RequestContext).(request.RequestContext)
	now := time.Now().UTC().Format(time.RFC3339Nano)
	audit := entity.NewAudit()

	audit.Operation.Tenant = rctx.Tenant
	audit.Operation.Action = action
	audit.Operation.Domain = domain
	audit.Operation.DomainId = domainId
	audit.Operation.CID = rctx.Cid
	audit.Operation.Date = now

	audit.User.Email = rctx.GetValueInCustomHeaders(request.HeaderXEmail)
	audit.User.Permission = rctx.Roles

	latitude, err := strconv.ParseFloat(rctx.GetValueInCustomHeaders(request.HeaderXLatitude), 64)
	if err != nil {
		audit.Localization.Latitude = nil
	} else {
		audit.Localization.Latitude = &latitude
	}

	longitude, err := strconv.ParseFloat(rctx.GetValueInCustomHeaders(request.HeaderXLongitude), 64)
	if err != nil {
		audit.Localization.Longitude = nil
	} else {
		audit.Localization.Longitude = &longitude
	}

	audit.UserAgent.Device = rctx.GetValueInCustomHeaders(request.HeaderUserAgent)
	audit.UserAgent.DeviceIp = utils.GetRequestIP(rctx.GetValueInCustomHeaders(request.HeaderXFowardedFor))

	audit.Http.Code = code
	audit.Http.Response = response
	audit.Http.Request = req

	return audit
}

