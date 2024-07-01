package otelx

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wakuwaku3/example-dapr-trace/server/lib/errorsx"
	"go.opentelemetry.io/otel/trace"
)

type (
	Trace struct {
		TraceID   string
		SpanID    string
		IsSampled bool
		IsRemote  bool
	}
)

func (t *Trace) String() string {
	bToI := func(b bool) int {
		if b {
			return 1
		}
		return 0
	}
	return fmt.Sprintf("%02d-%s-%s-%02d", bToI(t.IsRemote), t.TraceID, t.SpanID, bToI(t.IsSampled))
}

func GetTrace(ctx context.Context) *Trace {
	spanContext := trace.SpanContextFromContext(ctx)
	return &Trace{
		TraceID:   spanContext.TraceID().String(),
		SpanID:    spanContext.SpanID().String(),
		IsSampled: spanContext.TraceFlags().IsSampled(),
		IsRemote:  spanContext.IsRemote(),
	}
}

func GetTraceJSON(ctx context.Context) (json.RawMessage, error) {
	spanContext := trace.SpanContextFromContext(ctx)
	traceID, err := spanContext.TraceID().MarshalJSON()
	if err != nil {
		return nil, errorsx.Wrap(err)
	}
	return json.RawMessage(traceID), nil
}
