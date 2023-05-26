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
	enc.AddString("caller", err.Caller().String())
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
		err := enc.AddArray("context", encodeContextArray(data))
		if err != nil {
			enc.AddString("context", err.Error())
		}
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
