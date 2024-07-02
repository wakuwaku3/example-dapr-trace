package otelx

import (
	"log"

	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type (
	ZipkinTraceProviderFactory struct {
		URL         string
		Logger      *log.Logger
		ServiceName string
	}
)

func (z *ZipkinTraceProviderFactory) create() (*trace.TracerProvider, error) {
	exporter, err := zipkin.New(
		z.URL,
		zipkin.WithLogger(z.Logger),
	)
	if err != nil {
		return nil, err
	}

	batcher := trace.NewBatchSpanProcessor(exporter)

	return trace.NewTracerProvider(
		trace.WithSpanProcessor(batcher),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(z.ServiceName),
		)),
	), nil
}
