package otelx

import (
	"context"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	builder struct {
		name                  string
		traceProviderFactory  traceProviderFactory
		meterProviderFactory  meterProviderFactory
		loggerProviderFactory loggerProviderFactory
	}
	traceProviderFactory interface {
		create() (*trace.TracerProvider, error)
	}
	meterProviderFactory interface {
		create() (*metric.MeterProvider, error)
	}
	loggerProviderFactory interface {
		create() (*log.LoggerProvider, error)
	}
)

func NewBuilder(name string) *builder {
	return &builder{
		name: name,
	}
}

func (b *builder) WithTraceProviderFactory(factory traceProviderFactory) *builder {
	b.traceProviderFactory = factory
	return b
}

func (b *builder) WithMeterProviderFactory(factory meterProviderFactory) *builder {
	b.meterProviderFactory = factory
	return b
}

func (b *builder) WithLoggerProviderFactory(factory loggerProviderFactory) *builder {
	b.loggerProviderFactory = factory
	return b
}

func (b *builder) Build(ctx context.Context) (ShutdownFunc, error) {
	prop := b.newPropagator()
	otel.SetTextMapPropagator(prop)

	shutdownFuncs := &shutdownFuncs{}

	if b.traceProviderFactory != nil {
		tp, err := b.traceProviderFactory.create()
		if err != nil {
			return nil, shutdownFuncs.shutdownWithError(err, ctx)
		}
		shutdownFuncs.add(tp.Shutdown)
		otel.SetTracerProvider(tp)
		Provider.SetTracer(otel.Tracer(b.name))
	}

	if b.meterProviderFactory != nil {
		meterProvider, err := b.meterProviderFactory.create()
		if err != nil {
			return nil, shutdownFuncs.shutdownWithError(err, ctx)
		}
		shutdownFuncs.add(meterProvider.Shutdown)
		otel.SetMeterProvider(meterProvider)
		Provider.SetMeter(otel.Meter(b.name))
	}

	if b.loggerProviderFactory != nil {
		loggerProvider, err := b.loggerProviderFactory.create()
		if err != nil {
			return nil, shutdownFuncs.shutdownWithError(err, ctx)
		}
		shutdownFuncs.add(loggerProvider.Shutdown)
		global.SetLoggerProvider(loggerProvider)

		Provider.SetLogger(otelslog.NewLogger(b.name))
	}

	return shutdownFuncs.shutdown, nil
}

func (b *builder) newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}
