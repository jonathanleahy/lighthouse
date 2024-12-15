package network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/apierror"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/logger"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/request"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"runtime/debug"
	"strings"
	"time"
)

const ApiDurationExecMs = "api.duration_exec_ms"
const FailedResponseBodyMessage = "Failed to read response body"

type httpService struct {
	client           HttpClient
	auditService     AuditService
	sensitiveHeaders map[string]bool
}

type requestUrl struct {
	url             *url.URL
	gdprUrl         *url.URL
	pathParameters  []*PathParameter
	queryParameters []*QueryParameter
}

func NewHttpService(client HttpClient, auditService AuditService) HttpService {
	return &httpService{
		client:       client,
		auditService: auditService,
		sensitiveHeaders: map[string]bool{
			"x-email": true,
		},
	}
}

func (h *httpService) HttpRequest(ctx context.Context, req *Request) ([]byte, int, error) {
	rctx := ctx.Value(env.RequestContext).(request.RequestContext)

	err := req.validate()
	if err != nil {
		logger.Error("Failed to validate request", rctx.Cid, rctx.Tenant, logger.Fields{"error": err.Error(), "stacktrace": getStackTrace()})
		return nil, http.StatusInternalServerError, apierror.NewDomainError(ctx, http.StatusText(http.StatusInternalServerError), "Internal Server Error", nil)
	}

	reqUrl, err := buildRequestUrl(req.URL, req.PathParameters, req.QueryParameters)
	if err != nil {
		logger.Error("Failed to build url", rctx.Cid, rctx.Tenant, logger.Fields{"error": err.Error(), "stacktrace": getStackTrace()})
		return nil, http.StatusInternalServerError, apierror.NewDomainError(ctx, http.StatusText(http.StatusInternalServerError), "Internal Server Error", nil)
	}

	if req.Body == nil {
		return h.execute(ctx, req, reqUrl, nil)
	}

	if byteBody, ok := req.Body.([]byte); ok {
		return h.execute(ctx, req, reqUrl, byteBody)
	}

	byteBody, err := json.Marshal(req.Body)
	if err != nil {
		logger.Error("Failed to parse request body", rctx.Cid, rctx.Tenant, logger.Fields{"method": req.Method, "url": reqUrl.gdprUrl.String(), "status_code": http.StatusInternalServerError, "error": err.Error(), "stacktrace": getStackTrace()})
		return nil, http.StatusInternalServerError, apierror.NewDomainError(ctx, http.StatusText(http.StatusInternalServerError), "Internal Server Error", nil)
	}

	return h.execute(ctx, req, reqUrl, byteBody)
}

func (h *httpService) execute(ctx context.Context, req *Request, reqUrl *requestUrl, body []byte) (responseBody []byte, responseCode int, err error) {
	var (
		sp           oteltrace.Span
		httpRequest  *http.Request
		httpResponse *http.Response
		rctx         request.RequestContext
		headers      []*Header
		spanHeaders  []*Header
	)

	start := time.Now()

	if req.Headers == nil {
		headers = make([]*Header, 0)
	} else {
		headers = req.Headers
	}

	spanHeaders = make([]*Header, 0)

	if req.Span == nil {
		rctx = ctx.Value(env.RequestContext).(request.RequestContext)
	} else if !req.Span.Ignore {
		sp, ctx = tracer.GenerateChildSpanWithCtx(ctx, req.Span.Name)
		rctx = ctx.Value(env.RequestContext).(request.RequestContext)
	}

	var buffer *bytes.Buffer
	if body == nil {
		buffer = bytes.NewBuffer([]byte{})
	} else {
		buffer = bytes.NewBuffer(body)
	}

	if sp != nil {
		defer func() {
			h.finalizeSpan(rctx, req.Method, reqUrl, responseCode, spanHeaders, body, responseBody, err, sp, req.Span)
		}()
	}

	httpRequest, err = http.NewRequest(req.Method, reqUrl.url.String(), buffer)
	if err != nil {
		logger.Error(FailedResponseBodyMessage, rctx.Cid, rctx.Tenant, logger.Fields{ApiDurationExecMs: time.Since(start).Milliseconds(), "method": httpRequest.Method, "url": reqUrl.gdprUrl.String(), "status_code": responseCode, "error": err.Error(), "stacktrace": getStackTrace()})
		return nil, http.StatusInternalServerError, apierror.NewDomainError(ctx, http.StatusText(http.StatusInternalServerError), "Internal Server Error", nil)
	}

	spanHeaders = append(spanHeaders, h.buildRequestHeadersAndReturnSpanHeaders(ctx, body, httpRequest, rctx, headers)...)

	httpResponse, err = h.client.Do(httpRequest)
	if err != nil {
		logger.Error(FailedResponseBodyMessage, rctx.Cid, rctx.Tenant, logger.Fields{ApiDurationExecMs: time.Since(start).Milliseconds(), "method": httpRequest.Method, "url": reqUrl.gdprUrl.String(), "status_code": responseCode, "error": err.Error(), "stacktrace": getStackTrace()})
		return nil, http.StatusInternalServerError, apierror.NewDomainError(ctx, http.StatusText(http.StatusInternalServerError), "Internal Server Error", nil)
	}

	responseCode = httpResponse.StatusCode
	responseBody, err = io.ReadAll(httpResponse.Body)
	if err != nil {
		logger.Error(FailedResponseBodyMessage, rctx.Cid, rctx.Tenant, logger.Fields{ApiDurationExecMs: time.Since(start).Milliseconds(), "method": httpRequest.Method, "url": reqUrl.gdprUrl.String(), "status_code": responseCode, "error": err.Error(), "stacktrace": getStackTrace()})
		return nil, httpResponse.StatusCode, apierror.NewDomainError(ctx, http.StatusText(http.StatusInternalServerError), "Internal Server Error", nil)
	}

	defer h.closeBody(httpResponse, rctx, start, httpRequest, reqUrl, responseCode, err)

	responseBody, responseCode, err = h.parseAndHandleAudit(ctx, responseCode, rctx, httpRequest, reqUrl, body, httpResponse, responseBody, req.Audit, req.ResponsePtr)
	return responseBody, responseCode, err
}

func (h *httpService) buildRequestHeadersAndReturnSpanHeaders(ctx context.Context, body []byte, req *http.Request, rctx request.RequestContext, headers []*Header) []*Header {
	newHeadersMap := make(map[string]*Header)
	if body != nil {
		newHeadersMap[env.HeaderContentType] = &Header{Name: env.ContentTypeJSON, Value: env.ContentTypeJSON}
	}

	newHeadersMap[env.HeaderXCid] = &Header{Name: env.HeaderXCid, Value: rctx.Cid}
	newHeadersMap[env.HeaderXTenant] = &Header{Name: env.HeaderXTenant, Value: rctx.Tenant}
	newHeadersMap[env.HeaderXAccountID] = &Header{Name: env.HeaderXAccountID, Value: rctx.AccountID}

	for k, v := range rctx.CustomHeaders {
		if len(v) == 0 {
			continue
		}

		newHeadersMap[k] = &Header{Name: k, Value: v, Sensitive: h.sensitiveHeaders[strings.ToLower(k)]}
	}

	for _, he := range headers {
		if he.Value == nil {
			continue
		}

		newHeadersMap[he.Name] = &Header{Name: he.Name, Value: he.Value, Sensitive: he.Sensitive || h.sensitiveHeaders[strings.ToLower(he.Name)]}
	}

	spanHeaders := make([]*Header, 0)
	for k, v := range newHeadersMap {
		req.Header[k] = []string{fmt.Sprintf("%v", v.Value)}
		spanHeaders = append(spanHeaders, v)
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	return spanHeaders
}

func (h *httpService) parseAndHandleAudit(ctx context.Context, responseCode int, rctx request.RequestContext, req *http.Request, reqUrl *requestUrl, body []byte, resp *http.Response, responseBody []byte, audit *Audit, responsePtr interface{}) ([]byte, int, error) {
	if responseCode >= http.StatusInternalServerError {
		logger.Warn("External server error [5xx]", rctx.Cid, rctx.Tenant, logger.Fields{"method": req.Method, "url": reqUrl.gdprUrl.String(), "status_code": responseCode, "stacktrace": getStackTrace()})
		h.handleAudit(ctx, audit, body, false, responseBody, responseCode)
		return responseBody, responseCode, apierror.NewDomainError(ctx, http.StatusText(responseCode), http.StatusText(responseCode), nil)
	}

	if responseCode >= http.StatusBadRequest {
		if responseCode != http.StatusNotFound {
			logger.Warn("External client error [4xx]", rctx.Cid, rctx.Tenant, logger.Fields{"method": req.Method, "url": reqUrl.gdprUrl.String(), "status_code": responseCode, "stacktrace": getStackTrace()})
			h.handleAudit(ctx, audit, body, false, responseBody, responseCode)
		}
		return responseBody, responseCode, apierror.NewDomainError(ctx, http.StatusText(responseCode), http.StatusText(responseCode), nil)
	}

	if responsePtr != nil && len(responseBody) > 0 {
		err := json.Unmarshal(responseBody, responsePtr)
		if err != nil {
			logger.Error("Failed to parse response body", rctx.Cid, rctx.Tenant, logger.Fields{"method": req.Method, "url": reqUrl.gdprUrl.String(), "error": err, "stacktrace": getStackTrace()})
			h.handleAudit(ctx, audit, body, false, responseBody, responseCode)
			return responseBody, http.StatusInternalServerError, apierror.NewDomainError(ctx, http.StatusText(responseCode), "Internal Server Error", nil)
		}

		h.handleAudit(ctx, audit, body, true, responseBody, responseCode)
		return responseBody, responseCode, err
	}

	h.handleAudit(ctx, audit, body, false, responseBody, responseCode)
	return responseBody, resp.StatusCode, nil
}

func (h *httpService) handleAudit(ctx context.Context, audit *Audit, requestBody []byte, responseParsed bool, responseBody []byte, responseCode int) {
	if audit == nil || audit.Ignore {
		return
	}

	if requestBody == nil {
		requestBody = make([]byte, 0)
	}

	if responseBody == nil {
		responseBody = make([]byte, 0)
	}

	h.auditService.SendAuditEvent(ctx, audit.Domain, audit.Action, audit.DomainID(responseParsed, responseBody, responseCode), string(requestBody), string(responseBody), responseCode)
}

func (h *httpService) finalizeSpan(rctx request.RequestContext, method string, reqUrl *requestUrl, responseCode int, headers []*Header, body []byte, responseBody []byte, err error, sp oteltrace.Span, span *Span) {
	attributes := []attribute.KeyValue{
		attribute.String("Cid", rctx.Cid),
		attribute.String("OrgId", rctx.Tenant),
		attribute.String("AccountId", rctx.AccountID),
		attribute.String("http.method", method),
		attribute.String("http.url", strings.ReplaceAll(reqUrl.gdprUrl.String(), sensitivePlaceholderValueUrlEncoded, sensitivePlaceholderValue)),
		attribute.String("http.scheme", reqUrl.gdprUrl.Scheme),
		attribute.String("http.host", reqUrl.gdprUrl.Host),
		attribute.String("http.target", strings.ReplaceAll(reqUrl.gdprUrl.Path, sensitivePlaceholderValueUrlEncoded, sensitivePlaceholderValue)),
		attribute.String("http.server_name", reqUrl.gdprUrl.Hostname()),
		attribute.String("http.user_agent", env.GetEnvWithDefaultAsString(env.AppName, env.DefaultAppName)),
		attribute.Int("http.status_code", responseCode),
		attribute.String("http.status_text", http.StatusText(responseCode)),
	}

	for _, header := range headers {
		if header.Sensitive {
			attributes = append(attributes, h.buildSpanHeader(header.Name, sensitivePlaceholderValue))
		} else {
			attributes = append(attributes, h.buildSpanHeader(header.Name, header.Value))
		}
	}

	attributes = append(attributes, attribute.Int("http.request_content_length", len(body)))
	attributes = append(attributes, attribute.Int("http.response_content_length", len(responseBody)))

	if err == nil {
		sp.SetStatus(codes.Ok, "success")
	} else {
		sp.SetStatus(codes.Error, err.Error())
		attributes = append(attributes, attribute.String("error.trace", getStackTrace()))
	}

	attributes = append(attributes, span.Values...)
	sp.SetAttributes(attributes...)
	sp.End()
}

func (h *httpService) buildSpanHeader(key string, value interface{}) attribute.KeyValue {
	spanKey := fmt.Sprintf("http.request.header.%s", strings.ToLower(strings.ReplaceAll(key, "-", "_")))

	switch reflect.TypeOf(value).Kind() {
	case reflect.Array, reflect.Slice:
		return attribute.Array(spanKey, value)
	default:
		return attribute.Array(spanKey, []string{fmt.Sprintf("%v", value)})
	}
}

func (h *httpService) closeBody(resp *http.Response, rctx request.RequestContext, start time.Time, req *http.Request, reqUrl *requestUrl, responseCode int, err error) {
	bErr := resp.Body.Close()
	if bErr != nil {
		logger.Warn("Failed to close body", rctx.Cid, rctx.Tenant, logger.Fields{ApiDurationExecMs: time.Since(start).Milliseconds(), "method": req.Method, "url": reqUrl.gdprUrl.String(), "status_code": responseCode, "error": err.Error(), "stacktrace": getStackTrace()})
	}
}

func buildRequestUrl(urlFormat string, pathParameters []*PathParameter, queryParameters []*QueryParameter) (*requestUrl, error) {
	reqUrl, gdprUrl, err := formatUrlAndGdprUrl(urlFormat, pathParameters, queryParameters)
	if err != nil {
		return nil, err
	}

	if pathParameters == nil {
		pathParameters = make([]*PathParameter, 0)
	}

	if queryParameters == nil {
		queryParameters = make([]*QueryParameter, 0)
	}

	return &requestUrl{
		url:             reqUrl,
		gdprUrl:         gdprUrl,
		pathParameters:  pathParameters,
		queryParameters: queryParameters,
	}, nil
}

func getStackTrace() string {
	if env.GetEnvWithDefaultAsBoolean(env.StackTraceEnabled, env.DefaultStackTraceEnabled) {
		return string(debug.Stack())
	}
	return ""
}

