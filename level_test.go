package maleo

import "testing"

func TestLevel_String(t *testing.T) {
	tests := []struct {
		name string
		l    Level
		want string
	}{
		{
			name: "ErrorLevel",
			l:    ErrorLevel,
			want: "error",
		},
		{
			name: "Debug",
			l:    DebugLevel,
			want: "debug",
		},
		{
			name: "Warn",
			l:    WarnLevel,
			want: "warn",
		},
		{
			name: "Fatal",
			l:    FatalLevel,
			want: "fatal",
		},
		{
			name: "Panic",
			l:    PanicLevel,
			want: "panic",
		},
		{
			name: "Info",
			l:    InfoLevel,
			want: "info",
		},
		{
			name: "Unknown",
			l:    Level(100),
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
