package maleo

import (
	"context"
	"testing"
)

func TestNoopLogger_Log(t *testing.T) {
	type args struct {
		ctx   context.Context
		entry Entry
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestNoopLogger_Log",
			args: args{
				ctx:   context.Background(),
				entry: Global().NewEntry("TestNoopLogger_Log").Freeze(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			no := NoopLogger{}
			no.Log(tt.args.ctx, tt.args.entry)
		})
	}
}

func TestNoopLogger_LogError(t *testing.T) {
	type args struct {
		ctx context.Context
		err Error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestNoopLogger_LogError",
			args: args{
				ctx: context.Background(),
				err: Global().Bail("TestNoopLogger_LogError").Freeze(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			no := NoopLogger{}
			no.LogError(tt.args.ctx, tt.args.err)
		})
	}
}
