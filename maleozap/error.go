package maleozap

import (
	"encoding/json"

	"go.uber.org/zap/zapcore"

	"github.com/tigorlazuardi/maleo"
)

type Error struct {
	maleo.Error
}

type richJsonError struct{ error }

func (r richJsonError) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if r.error == nil {
		_ = enc.AddReflected("error", nil)
		return nil
	}
	enc.AddString("summary", r.error.Error())
	b, err := json.Marshal(r.error)
	if err != nil {
		enc.AddString("details", err.Error())
		return nil
	}
	switch {
	// ignore empty object and array
	case len(b) == 2 && b[0] == '{' && b[1] == '}':
	case len(b) == 2 && b[0] == '[' && b[1] == ']':
	default:
		return enc.AddReflected("details", json.RawMessage(b))
	}
	return nil
}

func (err Error) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("code", err.Code())
	enc.AddString("message", err.Message())
	_ = enc.AddObject("caller", Caller{err.Caller()})
	if key := err.Key(); key != "" {
		enc.AddString("key", key)
	}
	data := err.Context()
	if len(data) == 1 {
		err := encodeObject(enc, "context", data[0])
		if err != nil {
			enc.AddString("context", err.Error())
		}
	} else if len(data) > 1 {
		fields := toZapFields(data)
		_ = enc.AddObject("context", zapcore.ObjectMarshalerFunc(func(oe zapcore.ObjectEncoder) error {
			for _, field := range fields {
				field.AddTo(oe)
			}
			return nil
		}))
	}

	origin := err.Unwrap()
	if origin == nil {
		_ = enc.AddReflected("error", nil)
		return nil
	}

	if err := encodeObject(enc, "error", origin); err != nil {
		enc.AddString("error", err.Error())
	}
	return nil
}
