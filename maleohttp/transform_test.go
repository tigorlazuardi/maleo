package maleohttp

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
)

type errorJson struct {
	Message string `json:"message"`
}

func (e errorJson) MarshalJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (e errorJson) Error() string {
	return e.Message
}

func TestSimpleErrorTransformer_ErrorBodyTransform(t *testing.T) {
	type args struct {
		in0 context.Context
		err error
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "handled nil",
			args: args{
				in0: context.Background(),
				err: nil,
			},
			want: map[string]interface{}{"error": "[nil]"},
		},
		{
			name: "handled json marshaler implementor",
			args: args{
				in0: context.Background(),
				err: errorJson{Message: "test"},
			},
			want: map[string]interface{}{"error": errorJson{Message: "test"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := SimpleErrorTransformer{}
			if got := n.ErrorBodyTransform(tt.args.in0, tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorBodyTransform() = %v, want %v", got, tt.want)
			}
		})
	}
}
