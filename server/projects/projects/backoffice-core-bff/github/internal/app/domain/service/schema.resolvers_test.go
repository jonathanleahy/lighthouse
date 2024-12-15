package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pismo/backoffice-core-bff/internal/app/domain/entity"
	"github.com/pismo/backoffice-core-bff/internal/app/domain/service/mock"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/request"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolver_Query(t *testing.T) {
	resolver := &Resolver{}

	queryResolver := resolver.Query()

	assert.NotNil(t, queryResolver)
}

func TestResolver_Mutation(t *testing.T) {
	resolver := &Resolver{}

	mutationResolver := resolver.Mutation()

	assert.NotNil(t, mutationResolver)
}

func Test_queryResolver_Health(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock.NewMockHealthService(ctrl)
	s.EXPECT().GetMessage().Return(&entity.Health{}, nil)

	r := &queryResolver{
		Resolver: &Resolver{
			HealthService: s,
		},
	}

	_, err := r.Health(nil)

	assert.Nil(t, err)
}

func Test_queryResolver_Dispute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{})

	d := mock.NewMockDisputeService(ctrl)
	d.EXPECT().
		GetDisputeStatus(ctx, gomock.Any(), gomock.Any()).
		Return(make([]*entity.DisputeStatus, 0), nil)

	r := &queryResolver{
		Resolver: &Resolver{
			DisputeService: d,
		},
	}

	_, err := r.GetDisputeStatus(ctx, 1, nil)

	assert.Nil(t, err)
}

func Test_mutationResolver_FraudReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{})

	d := mock.NewMockDisputeService(ctrl)
	d.EXPECT().
		CreateFraudReport(ctx, gomock.Any(), gomock.Any()).
		Return(&entity.FraudReport{}, nil)

	r := &mutationResolver{
		Resolver: &Resolver{
			DisputeService: d,
		},
	}

	_, err := r.CreateFraudReport(ctx, 1, entity.FraudReportInput{})

	assert.Nil(t, err)
}

func Test_mutationResolver_UpdateDisputeStatusEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{})

	d := mock.NewMockDisputeService(ctrl)
	d.EXPECT().
		UpdateDisputeStatusEvent(ctx, gomock.Any(), gomock.Any()).
		Return(&entity.DisputeEvent{}, nil)

	r := &mutationResolver{
		Resolver: &Resolver{
			DisputeService: d,
		},
	}

	_, err := r.UpdateDisputeStatusEvent(ctx, 1, entity.DisputeEventInput{})

	assert.Nil(t, err)
}

