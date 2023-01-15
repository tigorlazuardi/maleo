package maleo

import (
	"errors"
	"testing"
)

func Test_multierror_Error(t *testing.T) {
	tests := []struct {
		name string
		m    multierror
		want string
	}{
		{
			name: "empty",
			m:    multierror{},
			want: "",
		},
		{
			name: "single",
			m:    multierror{errors.New("test")},
			want: "1. test",
		},
		{
			name: "multiple",
			m:    multierror{errors.New("test"), errors.New("test2")},
			want: "1. test; 2. test2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
