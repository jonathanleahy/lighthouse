package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sort"
	"time"

	otelcontrib "go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	TracerKey  = "crm_core_bff"
	TracerName = "crm_core_bff"
)

var (
	App string
)

func init() {
	App = env.GetEnvWithDefaultAsString(env.AppName, env.DefaultAppName)
}

// ConfigOpentelemetry middleware adds a `span` on tracer
func ConfigOpentelemetry(opts ...Option) echo.MiddlewareFunc {
	return start(opts...)
}

func start(opts ...Option) echo.MiddlewareFunc {
	cfg := config{}
	for _, opt := range opts {
		opt.apply(&cfg)
	}
	if cfg.TracerProvider == nil {
		cfg.TracerProvider = otel.GetTracerProvider()
	}
	tracer := cfg.TracerProvider.Tracer(
		TracerName,
		oteltrace.WithInstrumentationVersion(otelcontrib.SemVersion()),
	)
	if cfg.Propagators == nil {
		cfg.Propagators = otel.GetTextMapPropagator()
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			paths2ignore := []string{"/health", "/metrics"}
			if contains(paths2ignore, c.Path()) || c.Request().Header.Get("x-service") == "health-check" {
				err := next(c)
				if err != nil {
					c.Error(err)
				}
				return nil
			}

			start := time.Now()

			c.Set(TracerKey, tracer)
			request := c.Request()
			response := c.Response()

			savedCtx := request.Context()
			defer func() {
				request = request.WithContext(savedCtx)
				c.SetRequest(request)
			}()
			ctx := cfg.Propagators.Extract(savedCtx, propagation.HeaderCarrier(request.Header))

			commonLabels := []attribute.KeyValue{
				attribute.String("Cid", c.Request().Header.Get("x-cid")),
				attribute.String("OrgId", c.Request().Header.Get("x-tenant")),
			}
			opts := []oteltrace.SpanStartOption{
				oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", request)...),
				oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(request)...),
				oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(App, c.Path(), request)...),
				oteltrace.WithAttributes(semconv.HTTPClientAttributesFromHTTPRequest(request)...),
				oteltrace.WithAttributes(commonLabels...),
			}
			spanName := c.Path()
			if spanName == "" {
				spanName = fmt.Sprintf("HTTP %s route not found", request.Method)
			}

			ctx, span := tracer.Start(ctx, spanName, opts...)
			defer span.End()

			// pass the span through the request context
			c.SetRequest(request.WithContext(ctx))

			// serve the request to the next middleware
			err := next(c)
			if err != nil {
				span.SetAttributes(attribute.String("http.error", err.Error()))
				// invokes the registered HTTP error handler
				c.Error(err)
			}

			attrs := semconv.HTTPAttributesFromHTTPStatusCode(c.Response().Status)
			spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(c.Response().Status)
			span.SetAttributes(attrs...)
			span.SetStatus(spanStatus, spanMessage)

			latency := time.Since(start)
			fields := []zapcore.Field{
				zap.String("TraceId", span.SpanContext().TraceID().String()),
				zap.String("SpanId", span.SpanContext().SpanID().String()),
				zap.Any("Attributes", map[string]interface{}{
					"Cid":       request.Header.Get("x-cid"),
					"OrgId":     request.Header.Get("x-tenant"),
					"Status":    response.Status,
					"Host":      request.Host,
					"Size":      response.Size,
					"RemoteIp":  c.RealIP(),
					"UserAgent": request.UserAgent(),
					"Latency":   latency.String(),
					"Request":   fmt.Sprintf("%s %s", request.Method, request.RequestURI),
				},
				),
			}

			n := response.Status
			switch {
			case n >= 500:
				zap.L().With(zap.Error(err)).Error("Server error", fields...)
			case n >= 400:
				zap.L().With(zap.Error(err)).Warn("Client error", fields...)
			case n >= 300:
				zap.L().Info("Redirection", fields...)
			}

			return nil
		}
	}
}

func contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}

