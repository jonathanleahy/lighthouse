package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/pismo/backoffice-core-bff/internal/app/domain/entity"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/network"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/network/mock"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/request"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDisputeService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockNet := mock.NewMockHttpService(ctrl)

	service := NewDisputeService(mockNet)

	assert.NotNil(t, service)
}

func Test_disputeService_GetDisputeStatus(t *testing.T) {
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{})
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c := mock.NewMockHttpClientCustom(t, ctrl, 200, make([]*entity.DisputeStatus, 0), nil)
	h := network.NewHttpService(c, mock.NewMockAuditService(ctrl))
	d := NewDisputeService(h)
	param := 1

	got, err := d.GetDisputeStatus(ctx, 1, &param)

	assert.NoError(t, err)
	assert.NotNil(t, got)
}

func Test_disputeService_GetDisputeStatus_Error(t *testing.T) {
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{})
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c := mock.NewMockHttpClientCustom(t, ctrl, 500, nil, errors.New("error"))
	h := network.NewHttpService(c, mock.NewMockAuditService(ctrl))
	d := NewDisputeService(h)

	got, err := d.GetDisputeStatus(ctx, 1, nil)

	assert.Error(t, err)
	assert.Nil(t, got)
}

func Test_disputeService_CreateFraudReport(t *testing.T) {
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{})
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuditService := mock.NewMockAuditService(ctrl)
	mockAuditService.EXPECT().
		SendAuditEvent(gomock.Any(), gomock.Any(), gomock.Any(), "1", gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)
	c := mock.NewMockHttpClientCustom(t, ctrl, 200, &entity.FraudReport{}, nil)
	h := network.NewHttpService(c, mockAuditService)
	d := NewDisputeService(h)

	got, err := d.CreateFraudReport(ctx, 1, entity.FraudReportInput{})

	assert.NoError(t, err)
	assert.NotNil(t, got)
}

func Test_disputeService_CreateFraudReport_Error(t *testing.T) {
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{})
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c := mock.NewMockHttpClientCustom(t, ctrl, 500, nil, errors.New("error"))
	h := network.NewHttpService(c, mock.NewMockAuditService(ctrl))
	d := NewDisputeService(h)

	got, err := d.CreateFraudReport(ctx, 1, entity.FraudReportInput{})

	assert.Error(t, err)
	assert.Nil(t, got)
}

func Test_disputeService_UpdateDisputeStatusEvent(t *testing.T) {
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{})
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuditService := mock.NewMockAuditService(ctrl)
	mockAuditService.EXPECT().
		SendAuditEvent(gomock.Any(), gomock.Any(), gomock.Any(), "1", gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)
	c := mock.NewMockHttpClientCustom(t, ctrl, 200, &entity.DisputeEvent{}, nil)
	h := network.NewHttpService(c, mockAuditService)
	d := NewDisputeService(h)

	got, err := d.UpdateDisputeStatusEvent(ctx, 1, entity.DisputeEventInput{})

	assert.NoError(t, err)
	assert.NotNil(t, got)
}

func Test_disputeService_UpdateDisputeStatusEvent_Error(t *testing.T) {
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{})
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c := mock.NewMockHttpClientCustom(t, ctrl, 500, nil, errors.New("error"))
	h := network.NewHttpService(c, mock.NewMockAuditService(ctrl))
	d := NewDisputeService(h)

	got, err := d.UpdateDisputeStatusEvent(ctx, 1, entity.DisputeEventInput{})

	assert.Error(t, err)
	assert.Nil(t, got)
}

