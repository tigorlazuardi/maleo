package maleohttp

import (
	"bytes"
	"compress/gzip"
	"io"
	"reflect"
	"testing"
)

func TestGzipCompression_Compress(t *testing.T) {
	type fields struct {
		level int
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		want1   bool
		wantErr bool
	}{
		{
			name: "compress",
			fields: fields{
				level: gzip.BestCompression,
			},
			args: args{
				b: bytes.Repeat([]byte("hello world "), 200),
			},
			want: []byte{
				31, 139, 8, 0, 0, 0, 0, 0, 2, 255, 202, 72, 205, 201, 201, 87, 40, 207, 47, 202, 73, 81, 24, 101, 143,
				178, 71, 217, 163, 236, 81, 246, 40, 123, 148, 61, 202, 30, 101, 83, 206, 6, 4, 0, 0, 255, 255, 215, 43,
				80, 10, 96, 9, 0, 0,
			}, // compressed bytes of repeated "hello world " 200 times
			want1:   true,
			wantErr: false,
		},
		{
			name: "no compression on small data",
			fields: fields{
				level: gzip.BestCompression,
			},
			args: args{
				b: bytes.Repeat([]byte("hello world "), 10),
			},
			want:    bytes.Repeat([]byte("hello world "), 10),
			want1:   false,
			wantErr: false,
		},
		{
			name: "error on invalid compression level",
			fields: fields{
				level: -3,
			},
			args: args{
				b: bytes.Repeat([]byte("hello world "), 200),
			},
			want:    bytes.Repeat([]byte("hello world "), 200),
			want1:   false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := GzipCompression{
				level: tt.fields.level,
			}
			got, got1, err := g.Compress(tt.args.b)
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

func TestGzipCompression_ContentEncoding(t *testing.T) {
	type fields struct {
		level int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "output is always gzip",
			fields: fields{
				level: gzip.BestCompression,
			},
			want: "gzip",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := GzipCompression{
				level: tt.fields.level,
			}
			if got := g.ContentEncoding(); got != tt.want {
				t.Errorf("ContentEncoding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGzipCompression_StreamCompress(t *testing.T) {
	type fields struct {
		level int
	}
	type args struct {
		origin      io.Reader
		contentType string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
		wantOk bool
	}{
		{
			name: "compress",
			fields: fields{
				level: gzip.BestCompression,
			},
			args: args{
				origin:      bytes.NewReader(bytes.Repeat([]byte("hello world "), 200)),
				contentType: "text/plain; charset=utf-8",
			},
			want: []byte{
				31, 139, 8, 0, 0, 0, 0, 0, 2, 255, 202, 72, 205, 201, 201, 87, 40, 207, 47, 202, 73, 81, 24, 101, 143,
				178, 71, 217, 163, 236, 81, 246, 40, 123, 148, 61, 202, 30, 101, 83, 206, 6, 4, 0, 0, 255, 255, 215, 43,
				80, 10, 96, 9, 0, 0,
			}, // compressed bytes of repeated "hello world " 200 times
			wantOk: true,
		},
		{
			name: "no compress on non human readable content type",
			fields: fields{
				level: gzip.BestCompression,
			},
			args: args{
				origin:      bytes.NewReader(bytes.Repeat([]byte("this is an image"), 200)),
				contentType: "image/png",
			},
			want:   bytes.Repeat([]byte("this is an image"), 200),
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := GzipCompression{
				level: tt.fields.level,
			}
			out, ok := g.StreamCompress(tt.args.contentType, tt.args.origin)
			if ok != tt.wantOk {
				t.Errorf("StreamCompress() ok = %v, want %v", ok, tt.wantOk)
			}
			got, err := io.ReadAll(out)
			if err != nil {
				t.Errorf("StreamCompress() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StreamCompress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGzipCompression(t *testing.T) {
	tests := []struct {
		name string
		want *GzipCompression
	}{
		{
			name: "default compression level is gzip.DefaultCompression",
			want: &GzipCompression{
				level: gzip.DefaultCompression,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGzipCompression(); !reflect.DeepEqual(got.level, tt.want.level) {
				t.Errorf("NewGzipCompression() level = %v, want level %v", got.level, tt.want.level)
			}
		})
	}
}

func TestNewGzipCompressionWithLevel(t *testing.T) {
	type args struct {
		lvl int
	}
	tests := []struct {
		name string
		args args
		want *GzipCompression
	}{
		{
			name: "expected level output",
			args: args{
				lvl: gzip.BestCompression,
			},
			want: &GzipCompression{
				level: gzip.BestCompression,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGzipCompressionWithLevel(tt.args.lvl); !reflect.DeepEqual(got.level, tt.want.level) {
				t.Errorf("NewGzipCompressionWithLevel() = %v, want %v", got.level, tt.want.level)
			}
		})
	}
}
