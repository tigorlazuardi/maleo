package maleozap

import (
	"context"

	"github.com/tigorlazuardi/maleo"
	"go.uber.org/zap"
)

var _ maleo.Logger = (*Logger)(nil)

type TraceCapturer interface {
	CaptureTrace(ctx context.Context) []zap.Field
}

type TraceCapturerFunc func(ctx context.Context) []zap.Field

func (f TraceCapturerFunc) CaptureTrace(ctx context.Context) []zap.Field {
	return f(ctx)
}

type DisableField uint8

func (d DisableField) Has(flag DisableField) bool {
	return d&flag != 0
}

func (d *DisableField) Set(flag DisableField) {
	*d |= flag
}

const (
	DisableTime DisableField = 1 << iota
	DisableCaller
	DisableTrace
	DisableService
	DisableKey
	DisableCode
	DisableContext
	DisableError

	DisableNothing DisableField = 0
	DisableAll     DisableField = ^DisableNothing
)

type Logger struct {
	*zap.Logger
	tracer TraceCapturer
	flag   DisableField
}

func NewLogger(l *zap.Logger) *Logger {
	return &Logger{
		Logger: l,
		tracer: TraceCapturerFunc(func(ctx context.Context) []zap.Field { return nil }),
		flag:   DisableTime,
	}
}

func (l *Logger) SetTraceCapturer(capturer TraceCapturer) {
	l.tracer = capturer
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
		elements = append(elements, zap.Stringer("caller", entry.Caller()))
	}

	if !l.flag.Has(DisableContext) {
		data := entry.Context()
		if len(data) == 1 {
			elements = append(elements, toField("context", data[0]))
		} else if len(data) > 1 {
			elements = append(elements, zap.Array("context", encodeContextArray(entry.Context())))
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
		elements = append(elements, zap.Stringer("caller", err.Caller()))
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
			elements = append(elements, zap.Array("context", encodeContextArray(err.Context())))
		}
	}
	if !l.flag.Has(DisableError) {
		origin := err.Unwrap()
		elements = append(elements, toField("error", origin))
	}
	l.Logger.Log(translateLevel(err.Level()), err.Message(), elements...)
}
