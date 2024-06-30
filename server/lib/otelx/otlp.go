package otelx

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type (
	OtlpTraceProviderFactory struct {
		Context     context.Context
		ServiceName string
	}
	OtlpMeterProviderFactory struct {
		Context     context.Context
		ServiceName string
	}
	OtlpLoggerProviderFactory struct {
		Context     context.Context
		ServiceName string
	}
)

func (o *OtlpTraceProviderFactory) create() (*trace.TracerProvider, error) {
	traceExporter, err := otlptracegrpc.New(o.Context, otlptracegrpc.WithEndpoint("otel_collector:4317"), otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	batcher := trace.NewBatchSpanProcessor(traceExporter)

	traceProvider := trace.NewTracerProvider(
		trace.WithSpanProcessor(batcher),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(o.ServiceName),
		)),
	)
	return traceProvider, nil
}

func (o *OtlpMeterProviderFactory) create() (*metric.MeterProvider, error) {
	metricExporter, err := otlpmetricgrpc.New(o.Context, otlpmetricgrpc.WithEndpoint("otel_collector:4317"), otlpmetricgrpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
		metric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(o.ServiceName),
		)),
	)
	return meterProvider, nil
}

func (o *OtlpLoggerProviderFactory) create() (*log.LoggerProvider, error) {
	logExporter, err := otlploghttp.New(o.Context, otlploghttp.WithEndpoint("otel_collector:4318"), otlploghttp.WithInsecure())
	if err != nil {
		return nil, err
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
		log.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(o.ServiceName),
		)),
	)
	return loggerProvider, nil
}
