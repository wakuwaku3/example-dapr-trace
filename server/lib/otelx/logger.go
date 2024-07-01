package otelx

import (
	"context"

	"github.com/wakuwaku3/example-dapr-trace/server/lib/errorsx"
	"github.com/wakuwaku3/example-dapr-trace/server/lib/logx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type (
	loggerWrapper struct{}
)

func NewLogger() *loggerWrapper {
	return &loggerWrapper{}
}

var _ logx.Logger = &loggerWrapper{}

func (l *loggerWrapper) Info(ctx context.Context, message string, attributes ...logx.KeyValue) {
	attrs := append(convertAttributes(attributes...), attribute.KeyValue{
		Key:   "severity",
		Value: attribute.StringValue(string(logx.Info)),
	})
	if span := trace.SpanFromContext(ctx); span != nil {
		span.SetAttributes(attrs...)
		span.AddEvent(message, trace.WithAttributes(attrs...))
	}
	if logger := Provider.GetLogger(); logger != nil {
		logger.InfoContext(ctx, message, trace.WithAttributes(attrs...))
	}
}

func (l *loggerWrapper) Error(ctx context.Context, err error, attributes ...logx.KeyValue) {
	attrs := append(convertAttributes(attributes...), attribute.KeyValue{
		Key:   "severity",
		Value: attribute.StringValue(string(logx.Error)),
	}, attribute.KeyValue{
		Key:   "stacktrace",
		Value: attribute.StringValue(errorsx.StackTrace(err)),
	})

	if span := trace.SpanFromContext(ctx); span != nil {
		span.SetAttributes(attrs...)
		span.AddEvent(err.Error(), trace.WithAttributes(attrs...))
	}
	if logger := Provider.GetLogger(); logger != nil {
		logger.ErrorContext(ctx, err.Error(), trace.WithAttributes(attrs...))
	}
}

func convertAttributes(attributes ...logx.KeyValue) []attribute.KeyValue {
	var attrs []attribute.KeyValue
	for _, v := range attributes {
		attrs = append(attrs, attribute.KeyValue{
			Key:   attribute.Key(v.Key),
			Value: attribute.StringValue(v.Value),
		})
	}
	return attrs
}
