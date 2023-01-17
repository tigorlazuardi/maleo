package maleo

import (
	"context"
)

type KVString struct {
	Key, Value string
}

type TraceCapturer interface {
	// CaptureTrace captures the trace from the context.
	CaptureTrace(ctx context.Context) []KVString
}

type TraceCapturerFunc func(ctx context.Context) []KVString

func (f TraceCapturerFunc) CaptureTrace(ctx context.Context) []KVString {
	return f(ctx)
}

type NoopTraceCapturer struct{}

func (n NoopTraceCapturer) CaptureTrace(context.Context) []KVString {
	return nil
}
