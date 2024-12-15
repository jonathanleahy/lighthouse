package middleware

import (
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// config is used to configure the mux middleware.
type config struct {
	TracerProvider trace.TracerProvider
	Propagators    propagation.TextMapPropagator
}

// Option specifies instrumentation configuration options.
type Option interface {
	apply(*config)
}

