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
