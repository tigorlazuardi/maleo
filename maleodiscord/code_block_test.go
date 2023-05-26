package maleodiscord

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/kinbiko/jsonassert"
)

func TestJSONCodeBlockBuilder_Build(t *testing.T) {
	type args struct {
		value []any
	}
	tests := []struct {
		name    string
		args    args
		wantW   []string
		wantErr bool
	}{
		{
			name: "expected output - multiple values",
			args: args{
				value: []any{1},
			},
			wantW: []string{"```json", "1", "```"},
		},
		{
			name: "expected output - single value",
			args: args{
				value: []any{struct{}{}},
			},
			wantW: []string{"```json", "{}", "```"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			J := JSONCodeBlockBuilder{}
			w := &bytes.Buffer{}
			err := J.Build(w, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotW := strings.Fields(w.String())
			if !reflect.DeepEqual(gotW, tt.wantW) {
				t.Errorf("Build() gotW = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

type mockError struct {
	Message string
}

func (m mockError) Error() string {
	return m.Message
}

type mockJSONError struct{}

func (m mockJSONError) Error() string {
	return ""
}

func (m mockJSONError) MarshalJSON() ([]byte, error) {
	return []byte(`{"foo":"bar"}`), nil
}

type errorFunc func()

func (e errorFunc) Error() string {
	return ""
}

func TestJSONCodeBlockBuilder_BuildError(t *testing.T) {
	type args struct {
		e error
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name: "expected output",
			args: args{
				e: mockError{Message: "test error"},
			},
			wantErr: false,
			wantW:   `{"error":{"summary":"test error","details":{"Message": "test error"}}}`,
		},
		{
			name: "expected output - private details / minimum details",
			args: args{
				e: errors.New("bar"),
			},
			wantErr: false,
			wantW:   `{"error":"bar"}`,
		},
		{
			name: "expected output - nil error",
			args: args{
				e: nil,
			},
			wantErr: false,
			wantW:   `{"error":null}`,
		},
		{
			name: "expected output - json marshaler",
			args: args{
				e: mockJSONError{},
			},
			wantErr: false,
			wantW:   `{"foo":"bar"}`,
		},
		{
			name: "expected output - unserializable error",
			args: args{
				e: errorFunc(func() {}),
			},
			wantErr: true,
			wantW:   `{"foo":"bar"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			J := JSONCodeBlockBuilder{}
			w := &bytes.Buffer{}
			err := J.BuildError(w, tt.args.e)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("BuildError() unexpected error: %v", err)
				}
				return
			}
			gotW := strings.Split(w.String(), "\n")
			if len(gotW) < 3 {
				t.Fatal("expected number of lines to be three or bigger")
			}
			content := gotW[1 : len(gotW)-1]
			assert := jsonassert.New(t)
			assert.Assertf(strings.Join(content, ""), tt.wantW)
			if t.Failed() {
				fmt.Println(w.String())
			}
		})
	}
}

func Test_valueMarshaler_CodeBlockJSON(t *testing.T) {
	tests := []struct {
		name    string
		v       valueMarshaler
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.v.CodeBlockJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("CodeBlockJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CodeBlockJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}
