package otelx

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/wakuwaku3/example-dapr-trace/server/lib/errorsx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type (
	logger struct {
		ctx  context.Context
		span trace.Span
	}
	level string
)

const (
	Info  level = "info"
	Error level = "error"
)

func NewLogger(ctx context.Context) *logger {
	span := trace.SpanFromContext(ctx)
	return &logger{ctx: ctx, span: span}
}

func (l *logger) Info(message string, attributes ...attribute.KeyValue) {
	l.span.AddEvent(message, trace.WithAttributes(attribute.KeyValue{
		Key:   "severity",
		Value: attribute.StringValue(string(Info)),
	}), trace.WithAttributes(attributes...))
	log.Println(append([]interface{}{Info, message}, convertAttributes(attributes...)...)...)
}

func (l *logger) Error(err error, attributes ...attribute.KeyValue) {
	l.span.AddEvent(err.Error(), trace.WithAttributes(attribute.KeyValue{
		Key:   "severity",
		Value: attribute.StringValue(string(Error)),
	}), trace.WithAttributes(attribute.KeyValue{
		Key:   "stacktrace",
		Value: attribute.StringValue(errorsx.StackTrace(err)),
	}), trace.WithAttributes(attributes...))
	printError(err, convertAttributes(attributes...)...)
}

func convertAttributes(attributes ...attribute.KeyValue) []interface{} {
	a := []interface{}{}
	for _, v := range attributes {
		a = append(a, v)
	}
	return a
}

func printError(err error, args ...any) {
	log.Printf("%s %v\n", Error, err)
	converted := &errorsx.Body{}
	if ok := errors.As(err, &converted); ok {
		fmt.Println(converted.Args)
		fmt.Print(converted.Stack)
	}
	fmt.Println(args...)
}
