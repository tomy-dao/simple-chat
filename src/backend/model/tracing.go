package model

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// GetTracer returns the tracer from context
func GetTracer(reqCtx *RequestContext) trace.Tracer {
	if reqCtx == nil || reqCtx.ctx == nil {
		return otel.Tracer("default")
	}
	return otel.Tracer("local")
}

// GetSpanFromContext extracts span from context
func GetSpanFromContext(ctx context.Context) trace.Span {
	span := trace.SpanFromContext(ctx)
	return span
}

// StartSpan starts a new span from RequestContext
func (r *RequestContext) StartSpan(operation string) (context.Context, trace.Span) {
	tracer := GetTracer(r)
	ctx, span := tracer.Start(r.ctx, operation)
	r.ctx = ctx
	r.span = span
	return ctx, span
}

