package maleo

import (
	"errors"
	"fmt"
	"strings"
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

func Test_query_CollectErrors(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name       string
		args       args
		wantLength int
		test       func(t *testing.T, errs []Error)
	}{
		{
			name: "nil",
			args: args{
				err: nil,
			},
			wantLength: 0,
		},
		{
			name: "single",
			args: args{
				err: Global().Bail("foo").Freeze(),
			},
			wantLength: 1,
		},
		{
			name: "multiple and ordered",
			args: args{
				err: func() error {
					err := Global().Bail("foo").Freeze()
					return Global().Wrap(err).Message("bar").Freeze()
				}(),
			},
			wantLength: 2,
			test: func(t *testing.T, errs []Error) {
				if errs[0].Message() != "bar" {
					t.Errorf("expected first error to be 'bar', got %q", errs[0].Message())
				}
				if errs[1].Message() != "foo" {
					t.Errorf("expected second error to be 'foo', got %q", errs[1].Message())
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qu := query{}
			if got := qu.CollectErrors(tt.args.err); len(got) != tt.wantLength {
				t.Errorf("CollectErrors() = %v, want %v length", got, tt.wantLength)
			} else {
				if tt.test != nil {
					tt.test(t, got)
				}
			}
		})
	}
}

func Test_query_GetStack(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		test func(t *testing.T, stack []ErrorStack)
	}{
		{
			name: "nil",
			args: args{
				err: nil,
			},
			test: func(t *testing.T, stack []ErrorStack) {
				if len(stack) != 0 {
					t.Errorf("expected stack to be empty, got %v", stack)
				}
			},
		},
		{
			name: "single",
			args: args{
				err: Bail("foo").Freeze(),
			},
			test: func(t *testing.T, stack []ErrorStack) {
				if len(stack) != 1 {
					t.Errorf("expected stack to have 1 element, got %v", stack)
				}
				if stack[0].Error.Error() != "foo" {
					t.Errorf("expected first stack element to be 'foo', got %q", stack[0].Error.Error())
				}
				if !strings.HasSuffix(stack[0].Caller.File(), "query_test.go") {
					t.Errorf("expected first stack element to be from query_test.go, got %q", stack[0].Caller.File())
				}
			},
		},
		{
			name: "multiple - unsupported error type bottom",
			args: args{
				err: WrapFreeze(errors.New("foo"), "bar"),
			},
			test: func(t *testing.T, stack []ErrorStack) {
				if len(stack) != 1 {
					t.Errorf("expected stack to have 1 element, got %v", stack)
				}
			},
		},
		{
			name: "multiple - unsupported error type top-most",
			args: args{
				err: fmt.Errorf("foo: %w", BailFreeze("bar")),
			},
			test: func(t *testing.T, stack []ErrorStack) {
				if len(stack) != 1 {
					t.Errorf("expected stack to have 1 element, got %v", stack)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qu := query{}
			tt.test(t, qu.GetStack(tt.args.err))
		})
	}
}

func Test_query_TopError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want Error
	}{
		{
			name: "nil",
			args: args{
				err: nil,
			},
			want: nil,
		},
		{
			name: "multiple",
			args: args{
				err: func() error {
					err := Global().Bail("foo").Freeze()
					return Global().Wrap(err).Message("bar").Freeze()
				}(),
			},
			want: Global().Bail("bar").Freeze(),
		},
		{
			name: "multiple - wrapped by other",
			args: args{
				err: func() error {
					err := Global().Bail("foo").Freeze()
					return fmt.Errorf("bar: %w", err)
				}(),
			},
			want: Global().Bail("foo").Freeze(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qu := query{}
			got := qu.TopError(tt.args.err)
			if (tt.want == nil) != (got == nil) {
				t.Errorf("TopError() = %v, want %v", got, tt.want)
			}
			if got != nil && tt.want != nil {
				gotMsg := got.Message()
				wantMsg := tt.want.Message()
				if gotMsg != wantMsg {
					t.Errorf("TopError().Message() = %v, want %v", gotMsg, wantMsg)
				}
			}
		})
	}
}

func Test_query_GetMessage(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name        string
		args        args
		wantMessage string
	}{
		{
			name: "nil",
			args: args{
				err: nil,
			},
			wantMessage: "",
		},
		{
			name: "single",
			args: args{
				err: Global().Bail("foo").Freeze(),
			},
			wantMessage: "foo",
		},
		{
			name: "multiple",
			args: args{
				err: func() error {
					err := Global().Bail("foo").Freeze()
					return Global().Wrap(err).Message("bar").Freeze()
				}(),
			},
			wantMessage: "bar",
		},
		{
			name: "wrapped by other error",
			args: args{
				err: fmt.Errorf("foo: %w", Global().Bail("bar").Freeze()),
			},
			wantMessage: "bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qu := query{}
			if gotMessage := qu.GetMessage(tt.args.err); gotMessage != tt.wantMessage {
				t.Errorf("GetMessage() = %v, want %v", gotMessage, tt.wantMessage)
			}
		})
	}
}

func Test_query_BottomError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want Error
	}{
		{
			name: "nil",
			args: args{
				err: nil,
			},
			want: nil,
		},
		{
			name: "multiple",
			args: args{
				err: func() error {
					err := Global().Bail("foo").Freeze()
					return Global().Wrap(err).Message("bar").Freeze()
				}(),
			},
			want: Global().Bail("foo").Freeze(),
		},
		{
			name: "multiple - wrapped by other",
			args: args{
				err: func() error {
					err := Global().Bail("foo").Freeze()
					return fmt.Errorf("bar: %w", err)
				}(),
			},
			want: Global().Bail("foo").Freeze(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qu := query{}
			got := qu.BottomError(tt.args.err)
			if (tt.want == nil) != (got == nil) {
				t.Errorf("BottomError() = %v, want %v", got, tt.want)
			}
			if got != nil && tt.want != nil {
				gotMsg := got.Message()
				wantMsg := tt.want.Message()
				if gotMsg != wantMsg {
					t.Errorf("BottomError().Message() = %v, want %v", gotMsg, wantMsg)
				}
			}
		})
	}
}

func Test_query_Cause(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "actual nil",
			args: args{
				err: nil,
			},
			want: nil,
		},
		{
			name: "nil - under",
			args: args{
				err: WrapFreeze(nil, "foo"),
			},
			want: ErrNil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qu := query{}
			got := qu.Cause(tt.args.err)
			if (tt.want == nil) != (got == nil) {
				t.Errorf("Cause() = %v, want %v", got, tt.want)
			}
		})
	}
}
