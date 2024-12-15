package apierror

import (
	"context"
	"errors"
	"github.com/google/uuid"
	errMessage "github.com/pismo/backoffice-core-bff/internal/app/infrastructure/apierror/message"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/message"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/request"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var (
	caller   = "test"
	cid      = uuid.New().String()
	orgId    = uuid.New().String()
	code     = "testcode"
	msg      = "teste message"
	errorMsg = "testcode: teste message"

	errorMessage = errMessage.ErrorMessage{
		ErrorCode:   code,
		UserMessage: msg,
	}
)

func TestNewDebug(t *testing.T) {
	err := errors.New(msg)
	apierr := NewDebug(context.Background(), errorMessage, caller, cid, orgId, err)
	assert.Equal(t, errorMsg, apierr.Error())
	assert.Equal(t, message.Debug, apierr.GetSeverity())
	assert.Equal(t, 0, apierr.GetHttpStatus())
	assert.Equal(t, err, apierr.err)
}

func TestNewWarning(t *testing.T) {
	err := errors.New(msg)
	apierr := NewWarning(context.Background(), err, errorMessage, caller, cid, orgId, http.StatusBadRequest)
	assert.Equal(t, errorMsg, apierr.Error())
	assert.Equal(t, message.Warning, apierr.GetSeverity())
	assert.Equal(t, http.StatusBadRequest, apierr.GetHttpStatus())
	assert.Equal(t, err, apierr.err)
}

func TestNewError(t *testing.T) {
	err := errors.New(msg)
	apierr := NewError(context.Background(), err, errorMessage, caller, cid, orgId, http.StatusInternalServerError)
	assert.Equal(t, errorMsg, apierr.Error())
	assert.Equal(t, message.Error, apierr.GetSeverity())
	assert.Equal(t, http.StatusInternalServerError, apierr.GetHttpStatus())
	assert.Equal(t, err, apierr.err)
}

func TestForbiddenError(t *testing.T) {
	//args
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{})

	result := ForbiddenError(ctx)

	assert.NotNil(t, result)
}

func TestFowardError(t *testing.T) {
	//args
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{})

	err := errors.New(msg)
	result := FowardError(ctx, errorMessage, err, caller, http.StatusInternalServerError)

	assert.NotNil(t, result)
}

func TestFowardErrorWithoutResponseBody(t *testing.T) {
	//args
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{})

	result := FowardErrorWithResponseBody(ctx, errorMessage, nil, caller, http.StatusInternalServerError)

	assert.NotNil(t, result)
}

func TestFowardErrorWithResponseBody(t *testing.T) {
	//args
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{})
	responseBody := []byte{123, 34, 109, 101, 115, 115, 97, 103, 101, 34, 58, 34, 101, 114, 114, 111, 114, 34, 125}

	result := FowardErrorWithResponseBody(ctx, errorMessage, responseBody, caller, http.StatusInternalServerError)

	assert.NotNil(t, result)
}

func TestFowardErrorWithResponseBodyAndErrorUnmarshall(t *testing.T) {
	//args
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{})
	responseBody := []byte{0}

	result := FowardErrorWithResponseBody(ctx, errorMessage, responseBody, caller, http.StatusInternalServerError)

	assert.NotNil(t, result)
	assert.Equal(t, "ECPBFF0002: Error parsing utils body with message: invalid character '\\x00' looking for beginning of value", result.Error())
}

func TestNewDomainError(t *testing.T) {
	//args
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{})

	err := NewDomainError(ctx, "error_code", "message", nil)

	assert.NotNil(t, err)
}

func TestFowardErrorWithNilBody(t *testing.T) {
	//args
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{})

	result := FowardErrorWithResponseBody(ctx, errorMessage, nil, caller, 999)

	assert.NotNil(t, result)
}

func TestFowardWarnWithResponseBody(t *testing.T) {
	//args
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{})

	result := FowardWarnWithResponseBody(ctx, errorMessage, nil, caller, http.StatusNotFound)

	assert.NotNil(t, result)
}

func TestInternalServerError(t *testing.T) {
	//args
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{})

	result := InternalServerError(ctx, errors.New("error"))

	assert.NotNil(t, result)
}

