package otelx

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type (
	StdoutTraceProviderFactory struct {
		Context     context.Context
		ServiceName string
	}
	StdoutMeterProviderFactory struct {
		Context     context.Context
		ServiceName string
	}
	StdoutLoggerProviderFactory struct {
		Context     context.Context
		ServiceName string
	}
)

func (s *StdoutTraceProviderFactory) create() (*trace.TracerProvider, error) {
	traceExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(time.Second),
		),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(s.ServiceName),
		)),
	)
	return traceProvider, nil
}

func (s *StdoutMeterProviderFactory) create() (*metric.MeterProvider, error) {
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(3*time.Second)),
		),
		metric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(s.ServiceName),
		)),
	)
	return meterProvider, nil
}

func (s *StdoutLoggerProviderFactory) create() (*log.LoggerProvider, error) {
	logExporter, err := stdoutlog.New()
	if err != nil {
		return nil, err
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
		log.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(s.ServiceName),
		)),
	)
	return loggerProvider, nil
}
