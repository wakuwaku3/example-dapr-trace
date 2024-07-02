package fluentx

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/wakuwaku3/example-dapr-trace/lib/errorsx"
	"github.com/wakuwaku3/example-dapr-trace/lib/logx"
	"github.com/wakuwaku3/example-dapr-trace/lib/otelx"
)

type (
	logger struct {
		fluent *fluent.Fluent
	}
	Config = fluent.Config
)

var (
	_ logx.Logger = (*logger)(nil)
	_ io.Closer   = (*logger)(nil)
)

func NewLogger(config Config) (*logger, error) {
	fluent, err := fluent.New(config)
	if err != nil {
		return nil, err
	}
	return &logger{
		fluent: fluent,
	}, nil
}

func (l *logger) Close() error {
	return l.fluent.Close()
}

func (l *logger) post(tag string, data interface{}) {
	if err := l.fluent.Post(tag, data); err != nil {
		log.Println("fluentd との接続に失敗しました。", errorsx.Wrap(err), tag, data)
	}
}

func (l *logger) Info(ctx context.Context, message string, attributes ...logx.KeyValue) {
	tr := otelx.GetTrace(ctx)
	data := map[string]string{
		"trace":            tr.String(),
		"trace_id":         tr.TraceID,
		"span_id":          tr.SpanID,
		"trace_is_sampled": fmt.Sprintf("%t", tr.IsSampled),
		"trace_is_remote":  fmt.Sprintf("%t", tr.IsRemote),
		"severity":         string(logx.Info),
	}
	for _, v := range attributes {
		data[v.Key] = v.Value
	}
	l.post(message, data)
}

func (l *logger) Error(ctx context.Context, err error, attributes ...logx.KeyValue) {
	tr := otelx.GetTrace(ctx)
	data := map[string]string{
		"trace":            tr.String(),
		"trace_id":         tr.TraceID,
		"span_id":          tr.SpanID,
		"trace_is_sampled": fmt.Sprintf("%t", tr.IsSampled),
		"trace_is_remote":  fmt.Sprintf("%t", tr.IsRemote),
		"severity":         string(logx.Error),
		"stacktrace":       errorsx.StackTrace(err),
	}
	for _, v := range attributes {
		data[v.Key] = v.Value
	}
	l.post(fmt.Sprintf("ERROR: %s", err.Error()), data)
}
