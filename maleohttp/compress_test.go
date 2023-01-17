package maleohttp

import (
	"reflect"
	"testing"
)

func TestNoCompression_Compress(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		want1   bool
		wantErr bool
	}{
		{
			name: "no compression op",
			args: args{
				b: []byte("test"),
			},
			want:    []byte("test"),
			want1:   false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoCompression{}
			got, got1, err := n.Compress(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Compress() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Compress() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNoCompression_ContentEncoding(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "no compression returns no content encoding",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoCompression{}
			if got := n.ContentEncoding(); got != tt.want {
				t.Errorf("ContentEncoding() = %v, want %v", got, tt.want)
			}
		})
	}
}
