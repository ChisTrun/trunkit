package xcontext

import (
	"context"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func NewTraceContext(ctx context.Context, operation string) context.Context {
	if ctx == nil {
		ctx = context.TODO()
	}
	_, found := tracer.SpanFromContext(ctx)
	if found {
		return ctx
	}

	_, ctx = tracer.StartSpanFromContext(ctx, operation)
	return ctx
}
