package maleozap

type LoggerOption interface {
	Apply(*Logger)
}

type LoggerOptionFunc func(*Logger)

func (f LoggerOptionFunc) Apply(l *Logger) {
	f(l)
}

func WithDisableFieldFlag(fieldFlag DisabledField) LoggerOption {
	return LoggerOptionFunc(func(l *Logger) {
		l.flag = fieldFlag
	})
}

func WithTraceCapturer(capturer TraceCapturer) LoggerOption {
	return LoggerOptionFunc(func(l *Logger) {
		l.tracer = capturer
	})
}
