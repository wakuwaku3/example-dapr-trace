package httpx

import (
	"context"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func startSpan(r *http.Request) (context.Context, trace.Span, error) {
	tr := otel.GetTracerProvider().Tracer("dapr-diagnostics")
	spanContext, err := getParentSpanContext(r)
	if err != nil {
		return nil, nil, err
	}

	ctx, span := tr.Start(trace.ContextWithSpanContext(r.Context(), spanContext), r.Method+" "+r.RequestURI, trace.WithSpanKind(trace.SpanKindServer), trace.WithAttributes(semconv.HTTPMethodKey.String(r.Method), semconv.HTTPURLKey.String(r.RequestURI)))
	return ctx, span, nil
}

func getParentSpanContext(r *http.Request) (trace.SpanContext, error) {
	traceparent := r.Header.Get("Traceparent")
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
