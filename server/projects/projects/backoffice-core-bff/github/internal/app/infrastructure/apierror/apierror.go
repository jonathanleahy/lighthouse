package apierror

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/labstack/echo/v4"
	errMessage "github.com/pismo/backoffice-core-bff/internal/app/infrastructure/apierror/message"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/logger"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/message"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/request"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/tracer"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.opentelemetry.io/otel/codes"
	"net/http"
)

type (
	ApiError struct {
		ctx            context.Context
		httpStatusCode int
		err            error
		severity       int
		message        errMessage.ErrorMessage
	}
)

func (a ApiError) Error() string {
	return a.GetMessage().ErrorCode + ": " + a.GetMessage().UserMessage
}

func logError(caller, cid string, orgId string, err error) {
	intErr := err.(ApiError)

	switch intErr.GetSeverity() {
	case message.Warning:
		logger.Warn(err.Error(), cid, orgId, logger.Fields{"caller": caller, "code": intErr.GetMessage().ErrorCode})
	case message.Error:
		logger.Error("", cid, orgId, logger.Fields{"caller": caller, "code": intErr.GetMessage().ErrorCode, "internalMessage": err.Error()})
	case message.Debug:
		logger.Debug("Error on "+caller, cid, orgId, logger.Fields{"internalErrorMessage": err.Error(), "errorCode": intErr.GetMessage().ErrorCode})
	}
}

func newApiError(ctx context.Context, err error, errorMessage errMessage.ErrorMessage, caller, cid string, orgId string, httpStatusCode, severity int) (errReturned ApiError) {
	span, ctx := tracer.GenerateChildSpanWithCtx(ctx, caller)
	defer span.End()
	span.SetStatus(codes.Error, err.Error())

	apiErr := ApiError{ctx, httpStatusCode, err, severity, errorMessage}
	logError(caller, cid, orgId, apiErr)
	return apiErr
}

func NewDebug(ctx context.Context, errorMessage errMessage.ErrorMessage, caller, cid string, orgId string, err error) ApiError {
	return newApiError(ctx, err, errorMessage, caller, cid, orgId, 0, message.Debug)
}

func NewWarning(ctx context.Context, err error, errorMessage errMessage.ErrorMessage, caller, cid string, orgId string, httpStatusCode int) ApiError {
	return newApiError(ctx, err, errorMessage, caller, cid, orgId, httpStatusCode, message.Warning)
}

func NewError(ctx context.Context, err error, errorMessage errMessage.ErrorMessage, caller, cid string, orgId string, httpStatusCode int) ApiError {
	return newApiError(ctx, err, errorMessage, caller, cid, orgId, httpStatusCode, message.Error)
}

func NewDomainError(ctx context.Context, errorCode string, message string, params map[string]interface{}) *gqlerror.Error {
	rctx := ctx.Value(env.RequestContext).(request.RequestContext)
	if params == nil {
		params = make(map[string]interface{})
	}

	return &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: message,
		Extensions: map[string]interface{}{
			"error_code": errorCode,
			"params":     params,
			"tracking": map[string]interface{}{
				"tenant": rctx.Tenant,
				"email":  rctx.GetValueInCustomHeaders(request.HeaderXEmail),
				"cid":    rctx.Cid,
			},
		},
	}
}

func (a *ApiError) GetHttpStatus() int {
	return a.httpStatusCode
}

func (a *ApiError) GetMessage() errMessage.ErrorMessage {
	return a.message
}

func (a *ApiError) GetSeverity() int {
	return a.severity
}

func ForbiddenError(ctx context.Context) ApiError {
	rctx := ctx.Value(env.RequestContext).(request.RequestContext)
	msg := message.ErrForbidden
	err := errors.New(msg.UserMessage)

	return NewWarning(ctx, err, msg,
		"auth", rctx.Cid, rctx.Tenant, http.StatusForbidden)
}

func FowardError(ctx context.Context, errorMessage errMessage.ErrorMessage, err error, caller string, responseCode int) error {
	rctx := ctx.Value(env.RequestContext).(request.RequestContext)
	errorMessage.UserMessage = fmt.Sprintf(errorMessage.UserMessage, err.Error())
	return NewError(ctx, errors.New(errorMessage.UserMessage), errorMessage,
		caller, rctx.Cid, rctx.Tenant, responseCode)
}

func FowardWithResponseBody(ctx context.Context, severity int, errorMessage errMessage.ErrorMessage, responseBody []byte, caller string, responseCode int) error {
	rctx := ctx.Value(env.RequestContext).(request.RequestContext)
	var response echo.HTTPError

	if len(responseBody) == 0 {
		response.Message = http.StatusText(responseCode)
	} else {
		err := json.Unmarshal(responseBody, &response)
		if err != nil {
			return BodyParsingError(ctx, err)
		}
	}

	if response.Message == nil || response.Message == "" {
		response.Message = http.StatusText(responseCode)
	}

	errorMessage.UserMessage = fmt.Sprintf(errorMessage.UserMessage, response.Message)
	return newApiError(ctx, errors.New(errorMessage.UserMessage), errorMessage, caller, rctx.Cid, rctx.Tenant, responseCode, severity)
}

func FowardWarnWithResponseBody(ctx context.Context, errorMessage errMessage.ErrorMessage, responseBody []byte, caller string, responseCode int) error {
	return FowardWithResponseBody(ctx, message.Error, errorMessage, responseBody, caller, responseCode)
}

func FowardErrorWithResponseBody(ctx context.Context, errorMessage errMessage.ErrorMessage, responseBody []byte, caller string, responseCode int) error {
	return FowardWithResponseBody(ctx, message.Error, errorMessage, responseBody, caller, responseCode)
}

func BodyParsingError(ctx context.Context, err error) error {
	rctx := ctx.Value(env.RequestContext).(request.RequestContext)
	msg := message.ErrApiBodyParsing
	msg.UserMessage = fmt.Sprintf(msg.UserMessage, err.Error())
	return NewError(ctx, err, msg, message.CallerBodyParsingError, rctx.Cid, rctx.Tenant, http.StatusInternalServerError)
}

func InternalServerError(ctx context.Context, err error) error {
	rctx := ctx.Value(env.RequestContext).(request.RequestContext)
	return NewError(ctx, err, message.ErrInternalServerError, message.CallerApiError, rctx.Cid, rctx.Tenant, http.StatusInternalServerError)
}

