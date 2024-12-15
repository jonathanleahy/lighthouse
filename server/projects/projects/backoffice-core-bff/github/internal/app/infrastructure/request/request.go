package request

import (
	"strings"
)

const (
	HeaderXUtilsID        = "x-utils-id"
	HeaderXB3TraceID      = "x-b3-traceid"
	HeaderXB3SpanID       = "x-b3-spanid"
	HeaderXB3ParentspanID = "x-b3-parentspanid"
	HeaderXB3Sampled      = "x-b3-sampled"
	HeaderXOtSpanContext  = "x-ot-span-context"
	HeaderXVersion        = "x-version"
	HeaderXRoles          = "x-roles"
	HeaderXEmail          = "x-email"
	HeaderXLatitude       = "x-latitude"
	HeaderXLongitude      = "x-longitude"
	HeaderUserAgent       = "user-agent"
	HeaderXFowardedFor    = "x-forwarded-for"
	HeaderXRolesSeparator = ","
)

var (
	CustomHeaders = []string{
		HeaderXUtilsID,
		HeaderXB3TraceID,
		HeaderXB3SpanID,
		HeaderXB3ParentspanID,
		HeaderXB3Sampled,
		HeaderXOtSpanContext,
		HeaderXVersion,
		HeaderXRoles,
		HeaderXEmail,
		HeaderXLatitude,
		HeaderXLongitude,
		HeaderUserAgent,
		HeaderXFowardedFor,
	}
)

type (
	RequestContext struct {
		Tenant        string
		AccountID     string
		Cid           string
		AuditRequest  string
		Roles         []string
		CustomHeaders map[string]string
	}
)

func (rctx *RequestContext) HasRole(role string) bool {
	for _, r := range rctx.Roles {
		if strings.EqualFold(role, r) {
			return true
		}
	}

	return false
}

func (rctx *RequestContext) HasAnyRole(roles ...string) bool {
	for _, r := range roles {
		if rctx.HasRole(r) {
			return true
		}
	}

	return false
}

func (rctx *RequestContext) GetValueInCustomHeaders(headerName string) string {
	return rctx.CustomHeaders[headerName]
}

