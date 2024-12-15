package network

import (
	"errors"
	"net/http"
	"strings"
)

type Request struct {
	Method          string
	URL             string
	PathParameters  []*PathParameter
	QueryParameters []*QueryParameter
	Headers         []*Header
	Body            interface{}
	ResponsePtr     interface{}
	Audit           *Audit
	Span            *Span
}

func (r *Request) validate() error {
	if len(r.Method) == 0 {
		return errors.New(missingRequestMethod)
	}

	if len(r.URL) == 0 {
		return errors.New(missingRequestUrl)
	}

	if err := r.validateSpan(); err != nil {
		return err
	}

	if err := r.validateAudit(); err != nil {
		return err
	}

	return nil
}

func (r *Request) validateSpan() error {
	if r.Span == nil {
		return errors.New(missingRequestSpan)
	}

	return r.Span.validate()
}

func (r *Request) validateAudit() error {
	if r.Audit == nil {
		if r.isWriteOperation() {
			return errors.New(missingRequestAuditForWriteOperation)
		}

		return nil
	}

	return r.Audit.validate()
}

func (r *Request) isWriteOperation() bool {
	method := strings.ToUpper(strings.TrimSpace(r.Method))
	return method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch || method == http.MethodDelete
}

