// nolint:staticcheck
package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pismo/backoffice-core-bff/internal/app/domain/entity"
	"github.com/pismo/backoffice-core-bff/internal/app/domain/service/mock"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/request"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSendAuditEventSuccess(t *testing.T) {
	// args
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{})

	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	snsSender := mock.NewMockeventSender(ctrl)
	snsSender.EXPECT().
		SendEvent(ctx, "audit", entity.ACTION_CREATE, gomock.Any()).
		Return(nil).
		Times(1)

	// test
	svc := NewAuditService(snsSender)
	err := <-svc.SendAuditEvent(ctx, "audit", entity.ACTION_CREATE, "1", "", "", http.StatusOK)
	svc.WaitFinish()

	assert.Nil(t, err)
}

func TestSendAuditEventSuccessWithContext(t *testing.T) {
	// args
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{
		Tenant:       "tentant",
		Cid:          "cid",
		Roles:        []string{"owner", "owner"},
		AuditRequest: "{'id': 1}",
		CustomHeaders: map[string]string{
			request.HeaderXEmail:       "test@pismo.io",
			request.HeaderXFowardedFor: "127.0.0.1",
			request.HeaderUserAgent:    "Linux",
			request.HeaderXLatitude:    "-23.9975569",
			request.HeaderXLongitude:   "-46.2590205",
		},
	})

	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	snsSender := mock.NewMockeventSender(ctrl)
	snsSender.EXPECT().
		SendEvent(ctx, "audit", entity.ACTION_CREATE, gomock.Any()).
		Return(nil).
		Times(1)

	// test
	svc := NewAuditService(snsSender)
	err := <-svc.SendAuditEvent(ctx, "audit", entity.ACTION_CREATE, "1", "", "", http.StatusOK)
	svc.WaitFinish()

	assert.Nil(t, err)
}

