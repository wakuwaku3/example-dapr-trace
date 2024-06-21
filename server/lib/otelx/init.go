package otelx

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type (
	ZipkinOption struct {
		URL    string
		Logger *log.Logger
	}

	TraceProviderOption struct {
		ServiceName string
	}

	exporterFactory interface {
		create() (sdktrace.SpanExporter, error)
	}
)

func Init(traceProviderOption *TraceProviderOption, exporterFactory exporterFactory) (func(context.Context) error, error) {
	exporter, err := exporterFactory.create()
	if err != nil {
		return nil, err
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(traceProviderOption.ServiceName),
		)),
	)
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}

func (z *ZipkinOption) create() (sdktrace.SpanExporter, error) {
	return zipkin.New(
		z.URL,
		zipkin.WithLogger(z.Logger),
	)
}
