package httpx

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/wakuwaku3/example-dapr-trace/lib/errorsx"
	"github.com/wakuwaku3/example-dapr-trace/lib/otelx"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func startSpan(r *http.Request) (context.Context, trace.Span, error) {
	tr := otelx.Provider.GetTracer()
	spanContext, err := getParentSpanContext(r)
	if err != nil {
		return nil, nil, err
	}

	name := r.Method + " " + r.RequestURI
	ctx, span := tr.Start(trace.ContextWithSpanContext(r.Context(), spanContext), name, trace.WithSpanKind(trace.SpanKindServer), trace.WithAttributes(semconv.HTTPMethodKey.String(r.Method), semconv.HTTPURLKey.String(r.RequestURI)))
	if err := otelx.Count(ctx, name); err != nil {
		return nil, nil, errorsx.Wrap(err)
	}
	return ctx, span, nil
}

func getParentSpanContext(r *http.Request) (trace.SpanContext, error) {
	traceparent := GetTraceparent(r)
	spited := strings.Split(traceparent, "-")
	if len(spited) != 4 {
		fmt.Println("invalid traceparent", traceparent)
		return trace.SpanContext{}, nil
	}
	tid, err := trace.TraceIDFromHex(strings.Split(traceparent, "-")[1])
	if err != nil {
		return trace.SpanContext{}, err
	}
	sid, err := trace.SpanIDFromHex(strings.Split(traceparent, "-")[2])
	if err != nil {
		return trace.SpanContext{}, err
	}
	return trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: trace.FlagsSampled,
	}), nil
}
