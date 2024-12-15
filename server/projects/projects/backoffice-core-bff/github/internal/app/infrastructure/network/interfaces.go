package network

import (
	"context"
	"net/http"
)

type HttpService interface {
	HttpRequest(ctx context.Context, req *Request) ([]byte, int, error)
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AuditService interface {
	SendAuditEvent(ctx context.Context, domain, action, domainID string, request string, response string, responseCode int) chan error
}

