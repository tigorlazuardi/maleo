package maleodiscord

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/tigorlazuardi/maleo"
)

type read struct {
	io.Reader
	io.Closer
}

func TestDiscordMultipart(t *testing.T) {
	tests := []struct {
		name      string
		test      func(t *testing.T) callback
		wantCount int
		error     error
		message   string
		context   []any
		extraOpts []DiscordOption
	}{
		{
			name: "expected",
			test: func(t *testing.T) callback {
				return func(r *http.Request) {
					s := new(strings.Builder)
					r.Body = read{io.TeeReader(r.Body, s), r.Body}
					if r.Method != http.MethodPost {
						t.Errorf("expected method %s, got %s", http.MethodPost, r.Method)
					}
					if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
						t.Errorf("expected content type to be multipart/form-data, got %s", r.Header.Get("Content-Type"))
					}
					if err := r.ParseMultipartForm(32 << 20); err != nil {
						t.Fatalf("failed to parse multipart form: %v", err)
					}
					contains := func(h *multipart.FileHeader, value string) {
						file, err := h.Open()
						if err != nil {
							t.Fatalf("failed to open file: %v", err)
						}
						defer func(file multipart.File) {
							err := file.Close()
							if err != nil {
								t.Fatalf("failed to close file: %v", err)
							}
						}(file)
						b, err := io.ReadAll(file)
						if err != nil {
							t.Fatalf("failed to read file: %v", err)
						}
						if !strings.Contains(string(b), value) {
							t.Errorf("expected file to contain %s, got %s", value, string(b))
						}
					}
					for _, f := range r.MultipartForm.File {
						for _, h := range f {
							func(h *multipart.FileHeader) {
								switch {
								case strings.Contains(h.Filename, "_summary.md"):
									if h.Header.Get("Content-Type") != "text/markdown; charset=utf-8" {
										t.Errorf("expected content type to be 'text/markdown; charset=utf-8', got %s", h.Header.Get("Content-Type"))
									}
									if h.Size == 0 {
										t.Errorf("expected file size to be greater than 0, got %d", h.Size)
									}
									contains(h, strings.Repeat("a", 10000))
								case strings.Contains(h.Filename, "_error.json"):
									if h.Header.Get("Content-Type") != "application/json" {
										t.Errorf("expected content type to be 'application/json', got %s", h.Header.Get("Content-Type"))
									}
									if h.Size == 0 {
										t.Errorf("expected file size to be greater than 0, got %d", h.Size)
									}
									contains(h, strings.Repeat("a", 10000))
								case strings.Contains(h.Filename, "_error_stack.txt"):
									if h.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
										t.Errorf("expected content type to be 'text/plain; charset=utf-8', got %s", h.Header.Get("Content-Type"))
									}
									if h.Size == 0 {
										t.Errorf("expected file size to be greater than 0, got %d", h.Size)
									}
								case h.Filename == "payload_json":
									if h.Header.Get("Content-Type") != "application/json" {
										t.Errorf("expected content type to be 'application/json', got %s", h.Header.Get("Content-Type"))
									}
									if h.Size == 0 {
										t.Errorf("expected file size to be greater than 0, got %d", h.Size)
									}
									contains(h, strings.Repeat("a", 10000))
								default:
									t.Errorf("unexpected file %s", h.Filename)
								}
							}(h)
						}
					}
					if t.Failed() {
						fmt.Println(s.String())
					}
				}
			},
			wantCount: 1,
			error:     errors.New("test"),
			message:   strings.Repeat("a", 10000),
			context:   []any{maleo.F{"foo": "bar"}},
			extraOpts: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tow, _ := maleo.NewTestingMaleo()
			m := newMockClient(t, tt.test(t))
			defer m.Close()
			d := NewDiscordBot("", append([]DiscordOption{WithClient(m)}, tt.extraOpts...)...)
			tow.Register(d)
			if tt.error != nil {
				_ = tow.Wrap(tt.error).Message(tt.message).Context(tt.context...).Notify(context.Background())
			} else {
				_ = tow.NewEntry(tt.message).Context(tt.context...).Notify(context.Background())
			}
			err := tow.Wait(context.Background())
			if err != nil {
				t.Fatalf("maleo.Wait() error = %v", err)
			}
			m.Wait()
			if m.count != tt.wantCount {
				t.Errorf("m.count = %d, want = %d", m.count, tt.wantCount)
			}
		})
	}
}
