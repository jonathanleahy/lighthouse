package service

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"

	"github.com/pismo/backoffice-core-bff/internal/app/adapter/graphql/generated"
	"github.com/pismo/backoffice-core-bff/internal/app/domain/entity"
)

// CreateFraudReport is the resolver for the createFraudReport field.
func (r *mutationResolver) CreateFraudReport(ctx context.Context, disputeID int, input entity.FraudReportInput) (*entity.FraudReport, error) {
	return r.DisputeService.CreateFraudReport(ctx, disputeID, input)
}

// UpdateDisputeStatusEvent is the resolver for the updateDisputeStatusEvent field.
func (r *mutationResolver) UpdateDisputeStatusEvent(ctx context.Context, disputeID int, input entity.DisputeEventInput) (*entity.DisputeEvent, error) {
	return r.DisputeService.UpdateDisputeStatusEvent(ctx, disputeID, input)
}

// Health is the resolver for the health field.
func (r *queryResolver) Health(ctx context.Context) (*entity.Health, error) {
	return r.HealthService.GetMessage()
}

// GetDisputeStatus is the resolver for the getDisputeStatus field.
func (r *queryResolver) GetDisputeStatus(ctx context.Context, disputeID int, disputeInstallmentID *int) ([]*entity.DisputeStatus, error) {
	return r.DisputeService.GetDisputeStatus(ctx, disputeID, disputeInstallmentID)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

