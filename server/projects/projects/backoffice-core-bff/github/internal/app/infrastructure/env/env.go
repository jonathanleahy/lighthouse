package env

import (
	"os"
	"strconv"
)

type Key string

const (
	AppName                  = "APP_NAME"
	DefaultAppName           = "backoffice-core-bff"
	Version                  = "VERSION"
	DefaultVersion           = "1.0.0"
	Env                      = "ENV"
	DefaultEnv               = "dev"
	LogLevel                 = "LOG_LEVEL"
	DefaultLogLevel          = "info"
	StackTraceEnabled        = "STACK_TRACE_ENABLED"
	DefaultStackTraceEnabled = false
	// Server
	Host                   = "HOST"
	DefaultHost            = "0.0.0.0"
	Port                   = "PORT"
	DefaultPort            = "8080"
	HttpTimeout            = "HTTP_DEFAULT_TIMEOUT"
	DefaultHttpTimeout     = 60
	RequestContext     Key = "requestContext"
	// Opentelemetry
	OtlpUrl        = "OTEL_EXPORTER_OTLP_ENDPOINT_GO"
	DefaultOtlpUrl = "localhost:4317"
	// AWS
	AwsRegion              = "AWS_REGION"
	DefaultAwsRegion       = "sa-east-1"
	SnsConsoleAudit        = "SNS_CONSOLE_AUDIT"
	DefaultSnsConsoleAudit = "arn:aws:sns:sa-east-1:270036487593:console-audit"
	// Network
	HeaderXTenant     = "x-tenant"
	HeaderXCid        = "x-cid"
	HeaderXAccountID  = "x-account-id"
	HeaderContentType = "Content-Type"
	ContentTypeJSON   = "application/json;charset=utf8"
	// APIs
	DisputesApiURL        = "DISPUTES_API_URL"
	DefaultDisputesApiURL = "https://api-disputes-ext.pismolabs.io"
)

func GetEnvWithDefaultAsString(envKey string, defaultVal string) string {
	val := os.Getenv(envKey)
	if val == "" {
		return defaultVal
	}
	return val
}

func GetEnvWithDefaultAsInt(envKey string, defaultVal int) int {
	val := os.Getenv(envKey)
	if val == "" {
		return defaultVal
	}
	intValue, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return intValue
}

func GetEnvWithDefaultAsBoolean(envKey string, defaultVal bool) bool {
	val := os.Getenv(envKey)
	if val == "" {
		return defaultVal
	}
	boolValue, err := strconv.ParseBool(val)
	if err != nil {
		return defaultVal
	}
	return boolValue
}

