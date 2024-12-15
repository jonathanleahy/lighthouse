package service

import (
	"context"
	"github.com/pismo/backoffice-core-bff/internal/app/domain/entity"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/network"
	"net/http"
	"strconv"
)

type disputeService struct {
	http                        network.HttpService
	getDisputesStatusURL        string
	createFraudReportURL        string
	updateDisputeStatusEventURL string
}

func NewDisputeService(http network.HttpService) DisputeService {
	disputesUrl := env.GetEnvWithDefaultAsString(env.DisputesApiURL, env.DefaultDisputesApiURL)
	return &disputeService{
		http:                        http,
		getDisputesStatusURL:        disputesUrl + "/v1/disputes/%s/possible-status",
		createFraudReportURL:        disputesUrl + "/v1/disputes/%s/fraud-report",
		updateDisputeStatusEventURL: disputesUrl + "/v1/disputes/%s/event",
	}
}

func (d *disputeService) GetDisputeStatus(ctx context.Context, disputeId int, disputeInstallmentId *int) ([]*entity.DisputeStatus, error) {
	response := new([]*entity.DisputeStatus)
	_, _, err := d.http.HttpRequest(ctx, &network.Request{
		Method:          http.MethodGet,
		URL:             d.getDisputesStatusURL,
		PathParameters:  []*network.PathParameter{{Value: disputeId}},
		QueryParameters: []*network.QueryParameter{{Name: "disputeInstallmentId", Value: disputeInstallmentId}},
		ResponsePtr:     response,
		Span:            &network.Span{Name: "get_dispute_status"},
	})

	if err != nil {
		return nil, err
	}

	return *response, nil
}

func (d *disputeService) CreateFraudReport(ctx context.Context, disputeId int, input entity.FraudReportInput) (*entity.FraudReport, error) {
	response := new(entity.FraudReport)
	_, _, err := d.http.HttpRequest(ctx, &network.Request{
		Method:         http.MethodPost,
		URL:            d.createFraudReportURL,
		PathParameters: []*network.PathParameter{{Value: disputeId}},
		Body:           input,
		ResponsePtr:    response,
		Span:           &network.Span{Name: "create_fraud_report"},
		Audit: &network.Audit{
			Domain:   "dispute",
			Action:   "create",
			DomainID: network.PreProcessedDomainID(strconv.Itoa(disputeId)),
		},
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *disputeService) UpdateDisputeStatusEvent(ctx context.Context, disputeId int, input entity.DisputeEventInput) (*entity.DisputeEvent, error) {
	response := new(entity.DisputeEvent)
	_, _, err := d.http.HttpRequest(ctx, &network.Request{
		Method:         http.MethodPost,
		URL:            d.updateDisputeStatusEventURL,
		PathParameters: []*network.PathParameter{{Value: disputeId}},
		Body:           input,
		ResponsePtr:    response,
		Span:           &network.Span{Name: "update_dispute_status_event"},
		Audit: &network.Audit{
			Domain:   "dispute",
			Action:   "update",
			DomainID: network.PreProcessedDomainID(strconv.Itoa(disputeId)),
		},
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

