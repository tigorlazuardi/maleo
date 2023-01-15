package maleo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

type foo struct{}

func (f foo) String() string {
	return "foo"
}

type bar struct{}

func (b bar) Summary() string {
	return "bar"
}

func TestFields_Summary(t *testing.T) {
	tests := []struct {
		name string
		f    Fields
		want string
	}{
		{
			name: "empty",
			f:    Fields{},
			want: "",
		},
		{
			name: "single",
			f:    Fields{"a": 1},
			want: "a: 1",
		},
		{
			name: "multiple",
			f:    Fields{"a": 1, "b": 2},
			want: "a: 1\nb: 2",
		},
		{
			name: "nested",
			f:    Fields{"a": Fields{"b": 1}},
			want: "a: \n  b: 1",
		},
		{
			name: "nested multiple",
			f:    Fields{"a": Fields{"b": 1, "c": 2}},
			want: "a: \n  b: 1\n  c: 2",
		},
		{
			name: "nested multiple multiple",
			f:    Fields{"a": Fields{"b": 1, "c": 2}, "d": Fields{"e": 3, "f": 4}},
			want: "a: \n  b: 1\n  c: 2\nd: \n  e: 3\n  f: 4",
		},
		{
			name: "stringer",
			f:    Fields{"a": foo{}},
			want: "a: foo",
		},
		{
			name: "string like",
			f: Fields{
				"a": "foo",
				"b": []byte("bar"),
				"c": json.RawMessage(`{}`),
				"d": strings.Repeat("a", 100),
				"e": bytes.Repeat([]byte("a"), 100),
				"f": json.RawMessage(strings.Repeat("a", 100)),
			},
			want: "a: \"foo\"\nb: \"bar\"\nc: \"{}\"\nd: [...]\ne: [...]\nf: [...]",
		},
		{
			name: "null",
			f:    Fields{"a": nil},
			want: "a: null",
		},
		{
			name: "struct",
			f:    Fields{"a": struct{}{}},
			want: "a: [object]",
		},
		{
			name: "summary",
			f:    Fields{"a": bar{}},
			want: "a: bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Summary(); got != tt.want {
				t.Errorf("Summary() = %v, want \n%v", got, tt.want)
				fmt.Println(got)
			}
		})
	}
}
