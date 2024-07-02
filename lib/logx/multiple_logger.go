package logx

import "context"

type (
	multipleLogger struct {
		loggers []Logger
	}
)

var _ Logger = (*multipleLogger)(nil)

func NewMultipleLogger(loggers ...Logger) *multipleLogger {
	return &multipleLogger{
		loggers: loggers,
	}
}

func (m *multipleLogger) Register(loggers ...Logger) {
	m.loggers = append(m.loggers, loggers...)
}

func (m *multipleLogger) Info(ctx context.Context, message string, attributes ...KeyValue) {
	for _, logger := range m.loggers {
		logger.Info(ctx, message, attributes...)
	}
}

func (m *multipleLogger) Error(ctx context.Context, err error, attributes ...KeyValue) {
	for _, logger := range m.loggers {
		logger.Error(ctx, err, attributes...)
	}
}
