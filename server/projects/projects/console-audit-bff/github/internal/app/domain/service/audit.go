package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/pismo/console-audit-bff/internal/app/adapter/web/graphql/entity"
	"github.com/pismo/console-audit-bff/internal/app/adapter/web/presenter"
	"github.com/pismo/psm-sdk/psm/attribute/value"
	"github.com/pismo/psm-sdk/psm/domain/op"
	"github.com/pismo/psm-sdk/psm/domain/pagination"
	"github.com/pismo/psm-sdk/psm/network/http"
	"github.com/pismo/psm-sdk/psm/network/http/decoder"
	"github.com/pismo/psm-sdk/psm/schema"
)

type (
	AuditService interface {
		FindAuditByID(ctx context.Context, id int) (*presenter.Audit, error)
		SearchAudit(ctx context.Context, page *int, perPage *int, order *entity.Order, beginDate *string, endDate *string, user *string, actions []entity.AuditAction, domains []string, domainID *string, cid *string) (*presenter.SearchAudit, error)
		ListRoles(ctx context.Context) ([]*entity.Role, error)
		ListUserRoles(ctx context.Context, user *string, actions []entity.AuditAction, domains []string, domainID *string, cid *string) (*presenter.UserRoles, error)
		SearchRoles(ctx context.Context, email *string, feature *string) (presenter.Roles, error)
		CreateRole(ctx context.Context, role *entity.RoleInput) error
		UpdateRole(ctx context.Context, role *entity.RoleInput) error
		CreateFeature(ctx context.Context, feature *entity.FeatureInput, user *string, actions []entity.AuditAction, domains []string, domainID *string, cid *string) error
		UpdateFeature(ctx context.Context, feature *entity.FeatureInput, id int, user *string, actions []entity.AuditAction, domains []string, domainID *string, cid *string) error
		FindFeatureByID(ctx context.Context, id int) (*entity.Feature, error)
		SearchFeature(ctx context.Context, page *int, perPage *int, order *entity.Order, id *string, name *string, parentFeatureID *string, user *string, actions []entity.AuditAction, domains []string, domainID *string, cid *string) (*presenter.SearchFeature, error)
		DeleteFeature(ctx context.Context, id int) error
		DeleteRole(ctx context.Context, id int) error
		AttachRoleToUser(ctx context.Context, id int, email string) error
		DetachRoleToUser(ctx context.Context, id int, email string) error
		AttachRoleToFeature(ctx context.Context, id int, featureID int, action entity.ActionInput) error
		DetachRoleToFeature(ctx context.Context, id int, featureID int, action entity.ActionInput) error
		CreateEndpoint(ctx context.Context, endpoint *entity.EndpointInput, user *string, actions []entity.AuditAction, domains []string, domainID *string, cid *string) error
		UpdateEndpoint(ctx context.Context, endpoint *entity.EndpointInput, id int, user *string, actions []entity.AuditAction, domains []string, domainID *string, cid *string) error
		FindEndpointByID(ctx context.Context, id int) (*entity.Endpoint, error)
		SearchEndpoint(ctx context.Context, page *int, perPage *int, order *entity.Order, id *string, service_name *string, method *string, path *string, user *string, actions []entity.AuditAction, domains []string, domainID *string, cid *string) (*presenter.SearchEndpoint, error)
		DeleteEndpoint(ctx context.Context, id int) error
	}

	auditService struct {
		httpClient schema.HttpClient
		apiUrl     string
	}
)

func NewAuditService(httpClient schema.HttpClient, apiUrl string) AuditService {
	return &auditService{
		httpClient: httpClient,
		apiUrl:     apiUrl,
	}
}

func (s *auditService) FindAuditByID(ctx context.Context, id int) (*presenter.Audit, error) {
	response := s.httpClient.Execute(ctx, &http.Request{
		Method: "GET",
		Host:   s.apiUrl,
		Route:  http.Route("/v1/audit/:id").PathParam("id", value.Plain(id)),
		Span:   http.Span().Attr("id", value.Plain(id)),
	}, decoder.Json(new(presenter.Audit)))

	if response.Error() != nil {
		if response.StatusCode() == http.StatusNotFound {
			return nil, nil
		}

		return nil, op.NewError(http.StatusInternalServerError, response.Error().Error(), response.Error())
	}

	return response.DecodedBody().(*presenter.Audit), nil
}

func (s *auditService) SearchAudit(ctx context.Context, page *int, perPage *int, order *entity.Order, beginDate *string, endDate *string, user *string, actions []entity.AuditAction, domains []string, domainID *string, cid *string) (*presenter.SearchAudit, error) {
	response := s.httpClient.Execute(ctx, &http.Request{
		Method: "GET",
		Host:   s.apiUrl,
		Route: http.Route("/v1/audit").
			QueryParam("beginDate", value.Plain(beginDate)).
			QueryParam("endDate", value.Plain(endDate)).
			QueryParam("user", value.Plain(user)).
			QueryParam("action", value.Plain(actions)).
			QueryParam("domain", value.Plain(domains)).
			QueryParam("domainID", value.Plain(domainID)).
			QueryParam("cid", value.Plain(cid)).
			WithPagination(&pagination.Pagination{
				Page:    *page,
				PerPage: *perPage,
				Order:   []string{order.String()},
			}),
	}, decoder.Json(new(presenter.SearchAudit)))

	if response.Error() != nil {
		return nil, op.NewError(http.StatusInternalServerError, response.Error().Error(), response.Error())
	}

	return response.DecodedBody().(*presenter.SearchAudit), nil
}

func (s *auditService) ListRoles(ctx context.Context) ([]*entity.Role, error) {
	response := s.httpClient.Execute(ctx, &http.Request{
		Method: "GET",
		Host:   s.apiUrl,
		Route:  http.Route("/v1/role"),
	}, decoder.Json(new([]*entity.Role)))

	if response.Error() != nil {
		return nil, op.NewError(http.StatusInternalServerError, response.Error().Error(), response.Error())
	}

	return response.DecodedBody().([]*entity.Role), nil
}

func (s *auditService) ListUserRoles(ctx context.Context, user *string, actions []entity.AuditAction, domains []string, domainID *string, cid *string) (*presenter.UserRoles, error) {
	response := s.httpClient.Execute(ctx, &http.Request{
		Method: "GET",
		Host:   s.apiUrl,
		Route: http.Route("/v1/role/user").
			QueryParam("email", value.Plain(user)).
			QueryParam("domain", value.Plain(domains)).
			QueryParam("domainID", value.Plain(domainID)).
			QueryParam("cid", value.Plain(cid)),
	}, decoder.Json(new(presenter.UserRoles)))

	if response.Error() != nil {
		return nil, op.NewError(http.StatusInternalServerError, response.Error().Error(), response.Error())
	}

	return response.DecodedBody().(*presenter.UserRoles), nil
}

func (s *auditService) SearchRoles(ctx context.Context, email *string, feature *string) (presenter.Roles, error) {
	response := s.httpClient.Execute(ctx, &http.Request{
		Method: "GET",
		Host:   s.apiUrl,
		Route: http.Route("/v2/role").
			QueryParam("email", value.Plain(email)).
			QueryParam("feature", value.Plain(feature)),
	}, decoder.Json(new(presenter.Roles)))

	if response.Error() != nil {
		return nil, op.NewError(http.StatusInternalServerError, response.Error().Error(), response.Error())
	}

	return response.DecodedBody().(presenter.Roles), nil
}

func (s *auditService) CreateRole(ctx context.Context, role *entity.RoleInput) error {
	response := s.httpClient.Execute(ctx, &http.Request{
		Method: "POST",
		Host:   s.apiUrl,
		Route:  http.Route("/v1/role"),
		Body:   role,
		Span:   http.Span().Name("create_role"),
	}, nil)

	if response.StatusCode() == http.StatusBadRequest {
		return wrapBadRequestError(response)
	}

	return response.Error()
}

func (s *auditService) UpdateRole(ctx context.Context, role *entity.RoleInput) error {
	response := s.httpClient.Execute(ctx, &http.Request{
		Method: "PUT",
		Host:   s.apiUrl,
		Route:  http.Route("/v1/role/:id").PathParam("id", value.Plain(role.ID)),
		Body:   role,
		Span:   http.Span().Name("update_role"),
	}, nil)

	if response.StatusCode() == http.StatusBadRequest {
		return wrapBadRequestError(response)
	}

	return response.Error()
}

func (s *auditService) CreateFeature(ctx context.Context, feature *entity.FeatureInput, user *string, actions []entity.AuditAction, domains []string, domainID *string, cid *string) error {
	response := s.httpClient.Execute(ctx, &http.Request{
		Method: "POST",
		Host:   s.apiUrl,
		Route:  http.Route("/v1/feature"),
		Body:   feature,
	}, nil)

	if response.StatusCode() == http.StatusBadRequest {
		return wrapBadRequestError(response)
	}

	return response.Error()
}

func (s *auditService) UpdateFeature(ctx context.Context, feature *entity.FeatureInput, id int, user *string, actions []entity.AuditAction, domains []string, domainID *string, cid *string) error {
	response := s.httpClient.Execute(ctx, &http.Request{
		Method: "PUT",
		Host:   s.apiUrl,
		Route:  http.Route("/v1/feature/:id").PathParam("id", value.Plain(id)),
		Body:   feature,
	}, nil)

	if response.StatusCode() == http.StatusBadRequest {
		return wrapBadRequestError(response)
	}

	return response.Error()
}

func (s *auditService) FindFeatureByID(ctx context.Context, id int) (*entity.Feature, error) {
	response := s.httpClient.Execute(ctx, &http.Request{
		Method: "GET",
		Host:   s.apiUrl,
		Route:  http.Route("/v1/feature/:id").PathParam("id", value.Plain(id)),
		Span:   http.Span().Attr("id", value.Plain(id)),
	}, decoder.Json(new(entity.Feature)))

	if response.Error() != nil {
		if response.StatusCode() == http.StatusNotFound {
			return nil, nil
		}

		return nil, op.NewError(http.StatusInternalServerError, response.Error().Error(), response.Error())
	}

	return response.DecodedBody().(*entity.Feature), nil
}

func (s *auditService) SearchFeature(ctx context.Context, page *int, perPage *int, order *entity.Order, id *string, name *string, parentFeatureID *string, user *string, actions []entity.AuditAction, domains []string, domainID *string, cid *string) (*presenter.SearchFeature, error) {
	response := s.httpClient.Execute(ctx, &http.Request{
		Method: "GET",
		Host:   s.apiUrl,
		Route: http.Route("/v1/feature").
			QueryParam("id", value.Plain(id)).
			QueryParam("name", value.Plain(name)).
			QueryParam("parentFeatureID", value.Plain(parentFeatureID)).
			QueryParam("user", value.Plain(user)).
			QueryParam("action", value.Plain(actions)).
			QueryParam("domain", value.Plain(domains)).
			QueryParam("domainID", value.Plain(domainID)).
			QueryParam("cid", value.Plain(cid)).
			WithPagination(&pagination.Pagination{
				Page:    *page,
				PerPage: *perPage,
				Order:   []string{order.String()},
			}),
	}, decoder.Json(new(presenter.SearchFeature)))

	if response.Error() != nil {
		return nil, op.NewError(http.StatusInternalServerError, response.Error().Error(), response.Error())
	}

	return response.DecodedBody().(*presenter.SearchFeature), nil
}

func (s *auditService) DeleteRole(ctx context.Context, id int) error {
	response := s.httpClient.Execute(ctx, &http.Request{
		Method: "DELETE",
		Host:   s.apiUrl,
		Route:  http.Route("/v1/role/:id").PathParam("id", value.Plain(id)),
		Span:   http.Span().Name("delete_role"),
	}, nil)

	if response.StatusCode() == http.StatusBadRequest {
		return wrapBadRequestError(response)
	}

	return response.Error()
}

func (s *auditService) DeleteFeature(ctx context.Context, id int) error {
	response := s.httpClient.Execute(ctx, &http.Request{
		Method: "DELETE",
		Host:   s.apiUrl,
		Route:  http.Route("/v1/feature/:id").PathParam("id", value.Plain(id)),
	}, nil)

	if response.StatusCode() == http.StatusBadRequest {
		return wrapBadRequestError(response)
	}

	return response.Error()
}

func (s *auditService) AttachRoleToUser(ctx context.Context, id int, email string) error {
	response := s.httpClient.Execute(ctx, &http.Request{
		Method: "POST",
		Host:   s.apiUrl,
		Route:  http.Route("/v1/role/:id/attach/user/:email/").PathParam("id", value.Plain(id)).PathParam("email", value.Plain(email)),
		Span:   http.Span().Name("attach_role_to_user"),
	}, nil)

	if response.StatusCode() == http.StatusBadRequest {
		return wrapBadRequestError(response)
	}

	return response.Error()
}

func (s *auditService) DetachRoleToUser(ctx context.Context, id int, email string) error {
	response := s.httpClient.Execute(ctx, &http.Request{
		Method: "DELETE",
		Host:   s.apiUrl,
		Route:  http.Route("/v1/role/:id/detach/user/:email/").PathParam("id", value.Plain(id)).PathParam("email", value.Plain(email)),
		Span:   http.Span().Name("detach_role_to_user"),
	}, nil)

	if response.StatusCode() == http.StatusBadRequest {
		return wrapBadRequestError(response)
	}

	return response.Error()
}
func (s *auditService) AttachRoleToFeature(ctx context.Context, id int, featureID int, action entity.ActionInput) error {
	response := s.httpClient.Execute(ctx, &http.Request{
		Method: "POST",
		Host:   s.apiUrl,
		Route:  http.Route("/v1/role/:id/attach/feature/:feature_id/").PathParam("id", value.Plain(id)).PathParam("feature_id", value.Plain(featureID)),
		Span:   http.Span().Name("attach_role_to_feature"),
		Body:   action,
	}, nil)

	if response.StatusCode() == http.StatusBadRequest {
		return wrapBadRequestError(response)
	}

	return response.Error()
}

func (s *auditService) DetachRoleToFeature(ctx context.Context, id int, featureID int, actionRquest entity.ActionInput) error {
	response := s.httpClient.Execute(ctx, &http.Request{
		Method: "DELETE",
		Host:   s.apiUrl,
		Route: http.Route("/v1/role/:id/detach/feature/:feature_id/").
			PathParam("id", value.Plain(id)).
			PathParam("feature_id", value.Plain(featureID)),
		Span: http.Span().
			Name("detach_role_to_feature"),
		Body: actionRquest,
	}, nil)
	switch response.StatusCode() {
	case http.StatusBadRequest:
		return wrapBadRequestError(response)
	default:
		return response.Error()
	}
}

func (s *auditService) CreateEndpoint(ctx context.Context, endpoint *entity.EndpointInput, user *string, actions []entity.AuditAction, domains []string, domainID *string, cid *string) error {
	response := s.httpClient.Execute(ctx, &http.Request{
		Method: "POST",
		Host:   s.apiUrl,
		Route:  http.Route("/v1/endpoint"),
		Body:   endpoint,
	}, nil)

	if response.StatusCode() == http.StatusBadRequest {
		errorMessage := &presenter.ErrorMessage{}
		_ = json.Unmarshal(response.Body(), errorMessage)
		return errors.New(errorMessage.Message)
	}

	return response.Error()
}

func (s *auditService) UpdateEndpoint(ctx context.Context, endpoint *entity.EndpointInput, id int, user *string, actions []entity.AuditAction, domains []string, domainID *string, cid *string) error {
	response := s.httpClient.Execute(ctx, &http.Request{
		Method: "PUT",
		Host:   s.apiUrl,
		Route:  http.Route("/v1/endpoint/:id").PathParam("id", value.Plain(id)),
		Body:   endpoint,
	}, nil)

	if response.StatusCode() == http.StatusBadRequest {
		errorMessage := &presenter.ErrorMessage{}
		_ = json.Unmarshal(response.Body(), errorMessage)
		return errors.New(errorMessage.Message)
	}

	return response.Error()
}

func (s *auditService) FindEndpointByID(ctx context.Context, id int) (*entity.Endpoint, error) {
	response := s.httpClient.Execute(ctx, &http.Request{
		Method: "GET",
		Host:   s.apiUrl,
		Route:  http.Route("/v1/endpoint/:id").PathParam("id", value.Plain(id)),
		Span:   http.Span().Attr("id", value.Plain(id)),
	}, decoder.Json(new(entity.Endpoint)))

	if response.Error() != nil {
		if response.StatusCode() == http.StatusNotFound {
			return nil, nil
		}

		return nil, op.NewError(http.StatusInternalServerError, response.Error().Error(), response.Error())
	}

	return response.DecodedBody().(*entity.Endpoint), nil
}

func (s *auditService) SearchEndpoint(ctx context.Context, page *int, perPage *int, order *entity.Order, id *string, service_name *string, method *string, path *string, user *string, actions []entity.AuditAction, domains []string, domainID *string, cid *string) (*presenter.SearchEndpoint, error) {
	response := s.httpClient.Execute(ctx, &http.Request{
		Method: "GET",
		Host:   s.apiUrl,
		Route: http.Route("/v1/endpoint").
			QueryParam("id", value.Plain(id)).
			QueryParam("service_name", value.Plain(service_name)).
			QueryParam("method", value.Plain(method)).
			QueryParam("path", value.Plain(path)).
			QueryParam("user", value.Plain(user)).
			QueryParam("action", value.Plain(actions)).
			QueryParam("domain", value.Plain(domains)).
			QueryParam("domainID", value.Plain(domainID)).
			QueryParam("cid", value.Plain(cid)).
			WithPagination(&pagination.Pagination{
				Page:    *page,
				PerPage: *perPage,
				Order:   []string{order.String()},
			}),
	}, decoder.Json(new(presenter.SearchEndpoint)))

	if response.Error() != nil {
		return nil, op.NewError(http.StatusInternalServerError, response.Error().Error(), response.Error())
	}

	return response.DecodedBody().(*presenter.SearchEndpoint), nil
}

func (s *auditService) DeleteEndpoint(ctx context.Context, id int) error {
	response := s.httpClient.Execute(ctx, &http.Request{
		Method: "DELETE",
		Host:   s.apiUrl,
		Route:  http.Route("/v1/endpoint/:id").PathParam("id", value.Plain(id)),
	}, nil)

	if response.StatusCode() == http.StatusBadRequest {
		errorMessage := &presenter.ErrorMessage{}
		_ = json.Unmarshal(response.Body(), errorMessage)
		return errors.New(errorMessage.Message)
	}

	return response.Error()
}

func wrapBadRequestError(response schema.HttpResponse) error {
	errorMessage := &presenter.ErrorMessage{}
	_ = json.Unmarshal(response.Body(), errorMessage)
	return errors.New(errorMessage.Message)
}

