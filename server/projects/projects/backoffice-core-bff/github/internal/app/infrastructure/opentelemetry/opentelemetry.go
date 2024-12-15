package opentelemetry

import (
	"context"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/trace"
	"time"
)

var (
	OtlpUrl string
)

func init() {
	OtlpUrl = env.GetEnvWithDefaultAsString(env.OtlpUrl, env.DefaultOtlpUrl)
}

func InitTracer(ctx context.Context) *trace.TracerProvider {
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(OtlpUrl),
		otlptracegrpc.WithTimeout(2*time.Second),
	)
	if err != nil {
		logger.Panic("InitTracer failed", "", "", logger.Fields{"error": err})
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(traceExporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{}))
	return tp
}

func InitMetric(ctx context.Context) *controller.Controller {
	metricExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(OtlpUrl),
	)

	if err != nil {
		logger.Panic("failed to initialize stdoutmetric export pipeline", "", "", logger.Fields{"error": err})
	}

	pusher := controller.New(
		processor.New(
			simple.NewWithExactDistribution(),
			metricExporter,
		),
		controller.WithExporter(metricExporter),
		controller.WithCollectPeriod(2*time.Second),
	)

	err = pusher.Start(ctx)
	if err != nil {
		logger.Panic("failed to initialize metric controller", "", "", logger.Fields{"error": err})
	}
	return pusher
}

