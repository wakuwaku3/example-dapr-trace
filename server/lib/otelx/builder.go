package otelx

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	builder struct {
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

func NewBuilder() *builder {
	return &builder{
		traceProviderFactory:  &StdoutTraceProviderFactory{},
		meterProviderFactory:  &StdoutMeterProviderFactory{},
		loggerProviderFactory: &StdoutLoggerProviderFactory{},
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

	tp, err := b.traceProviderFactory.create()
	if err != nil {
		return nil, shutdownFuncs.shutdownWithError(err, ctx)
	}
	shutdownFuncs.add(tp.Shutdown)
	otel.SetTracerProvider(tp)

	meterProvider, err := b.meterProviderFactory.create()
	if err != nil {
		return nil, shutdownFuncs.shutdownWithError(err, ctx)
	}
	shutdownFuncs.add(meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	loggerProvider, err := b.loggerProviderFactory.create()
	if err != nil {
		return nil, shutdownFuncs.shutdownWithError(err, ctx)
	}
	shutdownFuncs.add(loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)

	return shutdownFuncs.shutdown, nil
}

func (b *builder) newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}
