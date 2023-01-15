package maleo

import "context"

type Logger interface {
	Log(ctx context.Context, entry Entry)
	LogError(ctx context.Context, err Error)
}

// NoopLogger that does nothing. The default logger that Maleo uses.
type NoopLogger struct{}

func (NoopLogger) Log(ctx context.Context, entry Entry)    {}
func (NoopLogger) LogError(ctx context.Context, err Error) {}
