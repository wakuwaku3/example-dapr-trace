package otelx

import (
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type (
	PrometheusMeterProviderFactory struct {
		ServiceName string
	}
)

func (p *PrometheusMeterProviderFactory) create() (*metric.MeterProvider, error) {
	metricExporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metricExporter),
		metric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(p.ServiceName),
		)),
	)
	return meterProvider, nil
}
