package logx

type (
	provider struct {
		multipleLogger *multipleLogger
	}
)

var Provider = &provider{
	multipleLogger: NewMultipleLogger(),
}

func (p *provider) Get() Logger {
	return p.multipleLogger
}

func (p *provider) Set(loggers ...Logger) {
	p.multipleLogger.Register(loggers...)
}
