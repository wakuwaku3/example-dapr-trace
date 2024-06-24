package otelx

import (
	"context"

	"github.com/wakuwaku3/example-dapr-trace/server/lib/errorsx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type (
	loggerWrapper struct {
		ctx  context.Context
		span trace.Span
	}
	level string
)

const (
	Info  level = "info"
	Error level = "error"
)

func NewLogger(ctx context.Context) *loggerWrapper {
	span := trace.SpanFromContext(ctx)
	return &loggerWrapper{ctx: ctx, span: span}
}

func (l *loggerWrapper) Info(message string, attributes ...attribute.KeyValue) {
	attrs := append(attributes, attribute.KeyValue{
		Key:   "severity",
		Value: attribute.StringValue(string(Info)),
	})
	l.span.SetAttributes(attrs...)
	l.span.AddEvent(message, trace.WithAttributes(attrs...))
	// log.Println(append([]interface{}{Info, message}, convertAttributes(attrs...)...)...)
	Logger.InfoContext(l.ctx, message, trace.WithAttributes(attrs...))
}

func (l *loggerWrapper) Error(err error, attributes ...attribute.KeyValue) {
	attrs := append(attributes, attribute.KeyValue{
		Key:   "severity",
		Value: attribute.StringValue(string(Error)),
	}, attribute.KeyValue{
		Key:   "stacktrace",
		Value: attribute.StringValue(errorsx.StackTrace(err)),
	})
	l.span.SetAttributes(attrs...)
	l.span.AddEvent(err.Error(), trace.WithAttributes(attrs...))
	// printError(err, convertAttributes(attrs...)...)
	Logger.ErrorContext(l.ctx, err.Error(), trace.WithAttributes(attrs...))
}

// func convertAttributes(attributes ...attribute.KeyValue) []interface{} {
// 	a := []interface{}{}
// 	for _, v := range attributes {
// 		a = append(a, v)
// 	}
// 	return a
// }

// func printError(err error, args ...any) {
// 	log.Printf("%s %v\n", Error, err)
// 	converted := &errorsx.Body{}
// 	if ok := errors.As(err, &converted); ok {
// 		fmt.Println(converted.Args)
// 		fmt.Print(converted.Stack)
// 	}
// 	fmt.Println(args...)
// }
