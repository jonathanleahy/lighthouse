package tracer

import (
	"context"
	"go.opentelemetry.io/otel"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func GenerateChildSpanWithCtx(parentCtx context.Context, spanName string) (oteltrace.Span, context.Context) {
	tracer := otel.Tracer("crm_core_bff")
	ctx, span := tracer.Start(
		parentCtx,
		spanName)
	return span, ctx
}

