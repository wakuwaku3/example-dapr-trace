package otelx

import (
	"context"

	"github.com/wakuwaku3/example-dapr-trace/server/lib/errorsx"
)

type (
	ShutdownFunc func(context.Context) error

	shutdownFuncs struct {
		funcs []ShutdownFunc
	}
)

func (s *shutdownFuncs) add(f ShutdownFunc) {
	s.funcs = append(s.funcs, f)
}

func (s *shutdownFuncs) shutdown(ctx context.Context) error {
	return s.shutdownWithError(nil, ctx)
}

func (s *shutdownFuncs) shutdownWithError(err error, ctx context.Context) error {
	builder := errorsx.NewMultipleErrorBuilder()
	builder.Append(err)
	for _, f := range s.funcs {
		if err := f(ctx); err != nil {
			builder.Append(errorsx.Wrap(err))
		}
	}
	return builder.Build()
}
