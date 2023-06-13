package maleozap

import (
	"go.uber.org/zap/zapcore"

	"github.com/tigorlazuardi/maleo"
)

type Caller struct {
	maleo.Caller
}

func (ca Caller) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("location", ca.Caller.String())
	enc.AddString("name", ca.Caller.ShortName())
	return nil
}
