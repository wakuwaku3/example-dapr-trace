package otelx

import (
	"log/slog"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type (
	provider struct {
		tracer trace.Tracer
		meter  metric.Meter
		logger *slog.Logger
	}
)

var (
	Provider = &provider{}
)

func (p *provider) GetTracer() trace.Tracer {
	return p.tracer
}

func (p *provider) SetTracer(tracer trace.Tracer) {
	p.tracer = tracer
}

func (p *provider) GetMeter() metric.Meter {
	return p.meter
}

func (p *provider) SetMeter(meter metric.Meter) {
	p.meter = meter
}

func (p *provider) GetLogger() *slog.Logger {
	return p.logger
}

func (p *provider) SetLogger(logger *slog.Logger) {
	p.logger = logger
}
