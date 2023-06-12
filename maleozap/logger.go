package maleozap

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/tigorlazuardi/maleo"
)

var _ maleo.Logger = (*Logger)(nil)

type TraceCapturer interface {
	CaptureTrace(ctx context.Context) []zap.Field
}

type TraceCapturerFunc func(ctx context.Context) []zap.Field

func (f TraceCapturerFunc) CaptureTrace(ctx context.Context) []zap.Field {
	return f(ctx)
}

type DisabledField uint8

func (d DisabledField) Has(flag DisabledField) bool {
	return d&flag != 0
}

func (d *DisabledField) Set(flag DisabledField) {
	*d |= flag
}

const (
	DisableTime DisabledField = 1 << iota
	DisableCaller
	DisableTrace
	DisableService
	DisableKey
	DisableCode
	DisableContext
	DisableError

	DisableNothing DisabledField = 0
	DisableAll     DisabledField = ^DisableNothing
)

type Logger struct {
	*zap.Logger
	tracer TraceCapturer
	flag   DisabledField
}

func New(l *zap.Logger) *Logger {
	l = l.WithOptions(zap.AddCallerSkip(4))
	return &Logger{
		Logger: l,
		tracer: TraceCapturerFunc(func(ctx context.Context) []zap.Field { return nil }),
		flag:   DisableTime,
	}
}

func (l *Logger) SetTraceCapturer(capturer TraceCapturer) {
	l.tracer = capturer
}

// SetDisabledFieldFlag sets the field from Maleo's Entry or Maleo's Error
// from being added to the zap log fields.
//
// e.g. SetDisableField(maleozap.DisableTime | maleozap.DisableCaller)
// will prevent data of Error.Time() or Error.Caller() from being printed.
func (l *Logger) SetDisabledFieldFlag(flag DisabledField) {
	l.flag = flag
}

func (l *Logger) Log(ctx context.Context, entry maleo.Entry) {
	elements := make([]zap.Field, 0, 16)
	if !l.flag.Has(DisableTime) {
		elements = append(elements, zap.Time("time", entry.Time()))
	}
	if !l.flag.Has(DisableTrace) {
		elements = append(elements, l.tracer.CaptureTrace(ctx)...)
	}
	if !l.flag.Has(DisableService) {
		elements = append(elements, zap.Object("service", service(entry.Service())))
	}
	if !l.flag.Has(DisableKey) {
		if key := entry.Key(); key != "" {
			elements = append(elements, zap.String("key", key))
		}
	}
	if !l.flag.Has(DisableCode) {
		code := entry.Code()
		if code != 0 {
			elements = append(elements, zap.Int("code", code))
		}
	}
	if !l.flag.Has(DisableCaller) {
		elements = append(elements, zap.Object("caller", Caller{entry.Caller()}))
	}

	if !l.flag.Has(DisableContext) {
		data := entry.Context()
		if len(data) == 1 {
			elements = append(elements, toField("context", data[0]))
		} else if len(data) > 1 {
			fields := toZapFields(data)
			obj := zapcore.ObjectMarshalerFunc(func(oe zapcore.ObjectEncoder) error {
				for _, f := range fields {
					f.AddTo(oe)
				}
				return nil
			})
			elements = append(elements, zap.Object("context", obj))
		}
	}

	l.Logger.Log(translateLevel(entry.Level()), entry.Message(), elements...)
}

func (l *Logger) LogError(ctx context.Context, err maleo.Error) {
	elements := make([]zap.Field, 0, 16)
	if !l.flag.Has(DisableTime) {
		elements = append(elements, zap.Time("time", err.Time()))
	}
	if !l.flag.Has(DisableTrace) {
		elements = append(elements, l.tracer.CaptureTrace(ctx)...)
	}
	if !l.flag.Has(DisableService) {
		elements = append(elements, zap.Object("service", service(err.Service())))
	}
	if !l.flag.Has(DisableCode) {
		elements = append(elements, zap.Int("code", err.Code()))
	}
	if !l.flag.Has(DisableCaller) {
		elements = append(elements, zap.Object("caller", Caller{err.Caller()}))
	}
	if !l.flag.Has(DisableKey) {
		if key := err.Key(); key != "" {
			elements = append(elements, zap.String("key", key))
		}
	}
	if !l.flag.Has(DisableContext) {
		data := err.Context()
		if len(data) == 1 {
			elements = append(elements, toField("context", data[0]))
		} else if len(data) > 1 {
			fields := toZapFields(data)
			obj := zapcore.ObjectMarshalerFunc(func(oe zapcore.ObjectEncoder) error {
				for _, f := range fields {
					f.AddTo(oe)
				}
				return nil
			})
			elements = append(elements, zap.Object("context", obj))
		}
	}
	if !l.flag.Has(DisableError) {
		origin := err.Unwrap()
		elements = append(elements, toField("error", origin))
	}

	l.Logger.Log(translateLevel(err.Level()), err.Message(), elements...)
}
