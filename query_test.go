package maleo

import (
	"fmt"
	"testing"
)

func Test_query_SearchCode(t *testing.T) {
	type args struct {
		err  error
		code int
	}
	tests := []struct {
		name      string
		args      args
		wantFound bool
	}{
		{
			name: "nil",
			args: args{
				err: nil,
			},
		},
		{
			name: "not found",
			args: args{
				err:  Global().Bail("foo").Code(400).Freeze(),
				code: 500,
			},
		},
		{
			name: "found",
			args: args{
				err:  Global().Bail("foo").Code(400).Freeze(),
				code: 400,
			},
			wantFound: true,
		},
		{
			name: "found (http)",
			args: args{
				err:  Global().Bail("foo").Code(7400).Freeze(),
				code: 400,
			},
			wantFound: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qu := query{}
			if gotFound := qu.SearchCode(tt.args.err, tt.args.code); (gotFound == nil) == tt.wantFound {
				t.Errorf("query.SearchCode() found = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func Test_query_SearchCodeHint(t *testing.T) {
	type args struct {
		err  error
		code int
	}
	tests := []struct {
		name      string
		args      args
		wantFound bool
	}{
		{
			name: "nil",
			args: args{
				err: nil,
			},
		},
		{
			name: "not found",
			args: args{
				err:  Global().Bail("foo").Code(400).Freeze(),
				code: 500,
			},
		},
		{
			name: "found",
			args: args{
				err:  Global().Bail("foo").Code(400).Freeze(),
				code: 400,
			},
			wantFound: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qu := query{}
			if gotFound := qu.SearchCodeHint(tt.args.err, tt.args.code); (gotFound == nil) == tt.wantFound {
				t.Errorf("query.SearchCode() found = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func Test_query_GetHTTPCode(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name: "nil",
			args: args{
				err: nil,
			},
			wantCode: 500,
		},
		{
			name: "500",
			args: args{
				err: Global().BailFreeze("foo"),
			},
			wantCode: 500,
		},
		{
			name: "404",
			args: args{
				err: Global().Bail("foo").Code(404).Freeze(),
			},
			wantCode: 404,
		},
		{
			name: "404 (wrapped by others)",
			args: args{
				err: func() error {
					err := Global().Bail("foo").Code(404).Freeze()
					return fmt.Errorf("wrap: %w", err)
				}(),
			},
			wantCode: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qu := query{}
			if gotCode := qu.GetHTTPCode(tt.args.err); gotCode != tt.wantCode {
				t.Errorf("GetHTTPCode() = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}

func Test_query_SearchHTTPCode(t *testing.T) {
	type args struct {
		err  error
		code int
	}
	tests := []struct {
		name      string
		args      args
		wantFound bool
	}{
		{
			name: "nil",
			args: args{
				err: nil,
			},
		},
		{
			name: "not found",
			args: args{
				err:  Global().Bail("foo").Code(400).Freeze(),
				code: 500,
			},
		},
		{
			name: "found",
			args: args{
				err:  Global().Bail("foo").Code(400).Freeze(),
				code: 400,
			},
			wantFound: true,
		},
		{
			name: "found (nested)",
			args: args{
				err: func() error {
					err := Global().Bail("foo").Code(400).Freeze()
					return Global().Wrap(err).Code(500).Freeze()
				}(),
				code: 400,
			},
			wantFound: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qu := query{}
			if got := qu.SearchHTTPCode(tt.args.err, tt.args.code); (got == nil) == tt.wantFound {
				t.Errorf("SearchHTTPCode() = %v, want %v", got, tt.wantFound)
			}
		})
	}
}
