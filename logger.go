package maleo

import "context"

// --8<-- [start:interface]

type Logger interface {
	Log(ctx context.Context, entry Entry)
	LogError(ctx context.Context, err Error)
}

// --8<-- [end:interface]

// NoopLogger that does nothing. The default logger that Maleo uses.
type NoopLogger struct{}

func (NoopLogger) Log(ctx context.Context, entry Entry)    {}
func (NoopLogger) LogError(ctx context.Context, err Error) {}
