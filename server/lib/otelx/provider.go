package otelx

import (
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
)

const name = "dapr-diagnostics"

var (
	Tracer = otel.Tracer(name)
	Meter  = otel.Meter(name)
	Logger = otelslog.NewLogger(name)
)
