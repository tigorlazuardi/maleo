package maleozap

import (
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/tigorlazuardi/maleo"
)

type Entry struct {
	maleo.Entry
}

func (e Entry) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("level", e.Level().String())
	enc.AddString("message", e.Message())
	if time.Since(e.Time()) > time.Second {
		enc.AddTime("time", e.Time())
	}
	if key := e.Key(); key != "" {
		enc.AddString("key", key)
	}
	enc.AddString("caller", e.Caller().String())
	if ctx := e.Context(); len(ctx) > 0 {
		if len(ctx) == 1 {
			err := encodeContextObject(enc, ctx[0])
			if err != nil {
				enc.AddString("context", err.Error())
			}
		} else {
			fields := toZapFields(ctx)
			_ = enc.AddObject("context", zapcore.ObjectMarshalerFunc(func(oe zapcore.ObjectEncoder) error {
				for _, field := range fields {
					field.AddTo(oe)
				}
				return nil
			}))
		}
	}
	return nil
}
