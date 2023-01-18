package maleozap

import (
	"github.com/tigorlazuardi/maleo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func encodeContextObject(oe zapcore.ObjectEncoder, data any) error {
	return encodeObject(oe, "context", data)
}

func encodeObject(oe zapcore.ObjectEncoder, key string, data any) error {
	switch value := data.(type) {
	case zapcore.ObjectMarshaler:
		return oe.AddObject(key, value)
	case zapcore.ArrayMarshaler:
		return oe.AddArray(key, value)
	case maleo.Fields:
		return oe.AddObject(key, fields(value))
	case maleo.Error:
		return oe.AddObject(key, Error{value})
	case maleo.Entry:
		return oe.AddObject(key, Entry{value})
	case error:
		return oe.AddObject(key, richJsonError{value})
	case map[string]any:
		return oe.AddObject(key, fields(value))
	default:
		return oe.AddReflected(key, value)
	}
}

func encodeContextArray(ctx []any) zapcore.ArrayMarshaler {
	return zapcore.ArrayMarshalerFunc(func(ae zapcore.ArrayEncoder) error {
		for _, value := range ctx {
			var err error
			switch value := value.(type) {
			case zapcore.ObjectMarshaler:
				err = ae.AppendObject(value)
			case zapcore.ArrayMarshaler:
				err = ae.AppendArray(value)
			case maleo.Fields:
				err = ae.AppendObject(fields(value))
			case maleo.Error:
				err = ae.AppendObject(Error{value})
			case maleo.Entry:
				err = ae.AppendObject(Entry{value})
			case error:
				err = ae.AppendObject(zapcore.ObjectMarshalerFunc(func(encoder zapcore.ObjectEncoder) error {
					return encoder.AddObject("error", richJsonError{value})
				}))
			case map[string]any:
				err = ae.AppendObject(fields(value))
			default:
				err = ae.AppendReflected(value)
			}
			if err != nil {
				ae.AppendString(err.Error())
			}
		}
		return nil
	})
}

func toField(key string, value any) zap.Field {
	switch value := value.(type) {
	case zapcore.ObjectMarshaler:
		return zap.Object(key, value)
	case zapcore.ArrayMarshaler:
		return zap.Array(key, value)
	case maleo.Fields:
		return zap.Object(key, fields(value))
	case maleo.Error:
		return zap.Object(key, Error{value})
	case maleo.Entry:
		return zap.Object(key, Entry{value})
	case error:
		return zap.Object(key, richJsonError{value})
	case map[string]any:
		return zap.Object(key, fields(value))
	default:
		return zap.Any(key, value)
	}
}

func translateLevel(lvl maleo.Level) zapcore.Level {
	switch lvl {
	case maleo.DebugLevel:
		return zapcore.DebugLevel
	case maleo.InfoLevel:
		return zapcore.InfoLevel
	case maleo.WarnLevel:
		return zapcore.WarnLevel
	case maleo.ErrorLevel:
		return zapcore.ErrorLevel
	case maleo.FatalLevel:
		return zapcore.FatalLevel
	case maleo.PanicLevel:
		return zapcore.PanicLevel
	default:
		return zapcore.InvalidLevel
	}
}
