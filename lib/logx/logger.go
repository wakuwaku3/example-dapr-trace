package logx

import "context"

type (
	KeyValue struct {
		Key   string
		Value string
	}
	Logger interface {
		Info(ctx context.Context, message string, attributes ...KeyValue)
		Error(ctx context.Context, err error, attributes ...KeyValue)
	}
)
