package maleohttp

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestNoopCloneBody_Bytes(t *testing.T) {
	tests := []struct {
		name string
		want []byte
	}{
		{
			name: "expected output",
			want: []byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopCloneBody{}
			if got := n.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoopCloneBody_CloneBytes(t *testing.T) {
	tests := []struct {
		name string
		want []byte
	}{
		{
			name: "expected output",
			want: []byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopCloneBody{}
			if got := n.CloneBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CloneBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoopCloneBody_Len(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "expected output",
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopCloneBody{}
			if got := n.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoopCloneBody_Read(t *testing.T) {
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		args    args
		wantN   int
		wantErr bool
	}{
		{
			name: "expected output",
			args: args{
				p: []byte{},
			},
			wantN:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			no := NoopCloneBody{}
			gotN, err := no.Read(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				if err != io.EOF {
					t.Errorf("expected error to be EOF")
				}
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Read() gotN = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestNoopCloneBody_Reader(t *testing.T) {
	tests := []struct {
		name string
		want BufferedReader
	}{
		{
			name: "expected output",
			want: NoopCloneBody{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			no := NoopCloneBody{}
			if got := no.Reader(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoopCloneBody_String(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "expected output",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopCloneBody{}
			if got := n.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoopCloneBody_Truncated(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "expected output",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NoopCloneBody{}
			if got := n.Truncated(); got != tt.want {
				t.Errorf("Truncated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bodyCloner_Bytes(t *testing.T) {
	type fields struct {
		ReadCloser io.ReadCloser
		clone      BufferedReadWriter
		limit      int
		cb         func(error)
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "expected output",
			fields: fields{
				ReadCloser: io.NopCloser(bytes.NewBufferString("test")),
				clone:      bytes.NewBufferString("test"),
				limit:      -1,
				cb:         nil,
			},
			want: []byte("test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &bodyCloner{
				ReadCloser: tt.fields.ReadCloser,
				clone:      tt.fields.clone,
				limit:      tt.fields.limit,
				cb:         tt.fields.cb,
			}
			if got := c.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bodyCloner_CloneBytes(t *testing.T) {
	type fields struct {
		ReadCloser io.ReadCloser
		clone      BufferedReadWriter
		limit      int
		cb         func(error)
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "expected output",
			fields: fields{
				ReadCloser: io.NopCloser(bytes.NewBufferString("test")),
				clone:      bytes.NewBufferString("test"),
				limit:      -1,
				cb:         nil,
			},
			want: []byte("test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &bodyCloner{
				ReadCloser: tt.fields.ReadCloser,
				clone:      tt.fields.clone,
				limit:      tt.fields.limit,
				cb:         tt.fields.cb,
			}
			if got := c.CloneBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CloneBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bodyCloner_Close(t *testing.T) {
	type fields struct {
		ReadCloser io.ReadCloser
		clone      BufferedReadWriter
		limit      int
		cb         func(error)
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "expected output",
			fields: fields{
				ReadCloser: io.NopCloser(bytes.NewBufferString("test")),
				clone:      bytes.NewBufferString("test"),
				limit:      -1,
				cb:         nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &bodyCloner{
				ReadCloser: tt.fields.ReadCloser,
				clone:      tt.fields.clone,
				limit:      tt.fields.limit,
				cb:         tt.fields.cb,
			}
			if err := c.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_bodyCloner_Len(t *testing.T) {
	type fields struct {
		ReadCloser io.ReadCloser
		clone      BufferedReadWriter
		limit      int
		cb         func(error)
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "expected output",
			fields: fields{
				ReadCloser: io.NopCloser(bytes.NewBufferString("test")),
				clone:      bytes.NewBufferString("test"),
				limit:      -1,
				cb:         nil,
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &bodyCloner{
				ReadCloser: tt.fields.ReadCloser,
				clone:      tt.fields.clone,
				limit:      tt.fields.limit,
				cb:         tt.fields.cb,
			}
			if got := c.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bodyCloner_Read(t *testing.T) {
	type fields struct {
		ReadCloser io.ReadCloser
		clone      BufferedReadWriter
		limit      int
		cb         func(error)
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantN      int
		wantErr    bool
		wantString string
	}{
		{
			name: "expected output",
			fields: fields{
				ReadCloser: io.NopCloser(bytes.NewBufferString("test")),
				clone:      bytes.NewBufferString("test"),
				limit:      -1,
				cb:         nil,
			},
			args: args{
				p: []byte("test"),
			},
			wantN:      4,
			wantErr:    false,
			wantString: "testtest",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &bodyCloner{
				ReadCloser: tt.fields.ReadCloser,
				clone:      tt.fields.clone,
				limit:      tt.fields.limit,
				cb:         tt.fields.cb,
			}
			gotN, err := c.Read(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Read() gotN = %v, want %v", gotN, tt.wantN)
			}
			if c.clone.String() != tt.wantString {
				t.Errorf("Read() got = %v, want %v", c.clone.String(), tt.wantString)
			}
		})
	}
}

func Test_bodyCloner_Reader(t *testing.T) {
	type fields struct {
		ReadCloser io.ReadCloser
		clone      BufferedReadWriter
		limit      int
		cb         func(error)
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "buffer should not be the same as inner reader",
			fields: fields{
				ReadCloser: io.NopCloser(bytes.NewBufferString("test")),
				clone:      bytes.NewBufferString("test"),
				limit:      -1,
				cb:         nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &bodyCloner{
				ReadCloser: tt.fields.ReadCloser,
				clone:      tt.fields.clone,
				limit:      tt.fields.limit,
				cb:         tt.fields.cb,
			}
			if got := c.Reader(); got == tt.fields.clone {
				t.Errorf("pointer of Reader() should be different. got = %v, origin = %v", got, tt.fields.clone)
			}
		})
	}
}

func Test_bodyCloner_Reset(t *testing.T) {
	type fields struct {
		ReadCloser io.ReadCloser
		clone      BufferedReadWriter
		limit      int
		cb         func(error)
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "expected output",
			fields: fields{
				ReadCloser: io.NopCloser(bytes.NewBufferString("test")),
				clone:      bytes.NewBufferString("test"),
				limit:      -1,
				cb:         nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &bodyCloner{
				ReadCloser: tt.fields.ReadCloser,
				clone:      tt.fields.clone,
				limit:      tt.fields.limit,
				cb:         tt.fields.cb,
			}
			c.Reset()
			if c.Len() != 0 {
				t.Errorf("Reset() should reset the buffer. got = %v, want %v", c.Len(), 0)
			}
		})
	}
}

func Test_bodyCloner_String(t *testing.T) {
	type fields struct {
		ReadCloser io.ReadCloser
		clone      BufferedReadWriter
		limit      int
		cb         func(error)
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "expected output",
			fields: fields{
				ReadCloser: io.NopCloser(bytes.NewBufferString("test")),
				clone:      bytes.NewBufferString("test"),
				limit:      -1,
				cb:         nil,
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &bodyCloner{
				ReadCloser: tt.fields.ReadCloser,
				clone:      tt.fields.clone,
				limit:      tt.fields.limit,
				cb:         tt.fields.cb,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bodyCloner_Truncated(t *testing.T) {
	type fields struct {
		ReadCloser io.ReadCloser
		clone      BufferedReadWriter
		limit      int
		cb         func(error)
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "expected output",
			fields: fields{
				ReadCloser: io.NopCloser(bytes.NewBufferString("test")),
				clone:      bytes.NewBufferString("test"),
				limit:      2,
				cb:         nil,
			},
			want: true,
		},
		{
			name: "expected output 2",
			fields: fields{
				ReadCloser: io.NopCloser(bytes.NewBufferString("test")),
				clone:      bytes.NewBufferString("test"),
				limit:      5,
				cb:         nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &bodyCloner{
				ReadCloser: tt.fields.ReadCloser,
				clone:      tt.fields.clone,
				limit:      tt.fields.limit,
				cb:         tt.fields.cb,
			}
			if got := c.Truncated(); got != tt.want {
				t.Errorf("Truncated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bodyCloner_Write(t *testing.T) {
	type fields struct {
		ReadCloser io.ReadCloser
		clone      BufferedReadWriter
		limit      int
		cb         func(error)
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantN      int
		wantErr    bool
		wantString string
	}{
		{
			name: "expected output",
			fields: fields{
				ReadCloser: io.NopCloser(bytes.NewBufferString("test")),
				clone:      bytes.NewBufferString("test"),
				limit:      -1,
				cb:         nil,
			},
			args: args{
				p: []byte("test"),
			},
			wantN:      4,
			wantErr:    false,
			wantString: "testtest",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &bodyCloner{
				ReadCloser: tt.fields.ReadCloser,
				clone:      tt.fields.clone,
				limit:      tt.fields.limit,
				cb:         tt.fields.cb,
			}
			gotN, err := c.Write(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Write() gotN = %v, want %v", gotN, tt.wantN)
			}
			if c.String() != tt.wantString {
				t.Errorf("Write() got = %v, want %v", c.String(), tt.wantString)
			}
		})
	}
}

func Test_bodyCloner_onClose(t *testing.T) {
	var called bool
	c := &bodyCloner{
		ReadCloser: io.NopCloser(bytes.NewBufferString("test")),
		clone:      bytes.NewBufferString("test"),
		limit:      -1,
	}
	c.onClose(func(err error) {
		called = true
	})
	err := c.Close()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatalf("expected callback to be called upon close")
	}
}

func Test_noopReadWriter_Reset(t *testing.T) {
	type fields struct {
		BufferedReader BufferedReader
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "expected output",
			fields: fields{
				BufferedReader: bytes.NewBufferString("test"),
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n2 := noopReadWriter{
				BufferedReader: tt.fields.BufferedReader,
			}
			n2.Reset()
			if n2.String() != tt.want {
				t.Errorf("n2.String() got = %v, want %v", n2.String(), tt.want)
			}
		})
	}
}

func Test_noopReadWriter_Write(t *testing.T) {
	type fields struct {
		BufferedReader BufferedReader
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantN   int
		wantErr bool
		want    string
	}{
		{
			name: "expected output",
			fields: fields{
				BufferedReader: bytes.NewBufferString("test"),
			},
			want:  "test",
			wantN: 3,
			args: args{
				p: []byte("foo"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n2 := noopReadWriter{
				BufferedReader: tt.fields.BufferedReader,
			}
			gotN, err := n2.Write(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Write() gotN = %v, want %v", gotN, tt.wantN)
			}
			if n2.String() != tt.want {
				t.Errorf("n2.String() got = %v, want %v", n2.String(), tt.want)
			}
		})
	}
}
