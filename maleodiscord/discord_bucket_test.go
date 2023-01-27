package maleodiscord

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/kinbiko/jsonassert"

	"github.com/tigorlazuardi/maleo"
	"github.com/tigorlazuardi/maleo/bucket"
)

type mockBucket map[string]any

func (m mockBucket) Upload(_ context.Context, files []bucket.File) []bucket.UploadResult {
	results := make([]bucket.UploadResult, len(files))
	for i, f := range files {
		url := "https://example.com/" + f.Filename()
		results[i] = bucket.UploadResult{
			File: f,
			URL:  "https://example.com/" + f.Filename(),
		}
		m[url] = results[i]
	}
	return results
}

func TestBucket(t *testing.T) {
	tests := []struct {
		name       string
		test       func(t *testing.T) callback
		testBucket func(t *testing.T, b mockBucket)
		wantCount  int
		error      error
		message    string
		context    []any
		extraOpts  []DiscordOption
	}{
		{
			name: "should upload file",
			test: func(t *testing.T) callback {
				return func(r *http.Request) {
					if r.Method != http.MethodPost {
						t.Errorf("want method %s, got %s", http.MethodPost, r.Method)
					}
					if r.Header.Get("Content-Type") != "application/json" {
						t.Errorf("want content type %s, got %s", "application/json", r.Header.Get("Content-Type"))
					}
					body, err := io.ReadAll(r.Body)
					if err != nil {
						t.Fatalf("failed to read body: %v", err)
					}
					j := jsonassert.New(t)
					j.Assertf(string(body), `
					{
					  "content": "@here an error has occurred from service **test** on type **test** on environment **test**",
					  "embeds": [
						{
						  "title": "Summary",
						  "type": "rich",
						  "description": "<<PRESENCE>>",
						  "color": 1606980
						},
						{
						  "title": "Error",
						  "type": "rich",
						  "description": "<<PRESENCE>>",
						  "color": 7405835
						},
						{
							"title": "Error Stack",
							"type": "rich",
							"description": "<<PRESENCE>>",
							"color": 6098454
						},
						{
							"title": "Metadata",
							"type": "rich",
							"description": "<<PRESENCE>>",
							"timestamp": "<<PRESENCE>>",
							"color": 6576731,
							"fields": [
									{
										"name": "Service",
										"value": "test",
										"inline": true
									},
									{
										"name": "Type",
										"value": "test",
										"inline": true
									},
									{
										"name": "Environment",
										"value": "test",
										"inline": true
									},
									{
										"name": "Thread ID",
										"value": "<<PRESENCE>>",
										"inline": true
									},
									{
										"name": "Message Iteration",
										"value": "1",
										"inline": true
									},
									{
										"name": "Next Possible Earliest Repeat",
										"value": "<<PRESENCE>>"
									}
								]
							},
							{
								"title": "Attachments",
								"type": "rich",
								"color": 1606980,
								"fields": "<<PRESENCE>>"
							}
						]
					}`)
					if t.Failed() {
						out := new(bytes.Buffer)
						_ = json.Indent(out, body, "", "  ")
						fmt.Println(out.String())
					}
				}
			},
			testBucket: func(t *testing.T, b mockBucket) {
				if len(b) != 2 {
					t.Errorf("want 2, got %d", len(b))
					if t.Failed() {
						out := new(bytes.Buffer)
						enc := json.NewEncoder(out)
						enc.SetEscapeHTML(true)
						_ = enc.Encode(b)
						fmt.Printf(out.String())
					}
				}
			},
			wantCount: 1,
			error:     errors.New(strings.Repeat("foo", 3000)),
			message:   "",
			context:   nil,
			extraOpts: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mBucket := mockBucket{}
			mal, _ := maleo.NewTestingMaleo()
			m := newMockClient(t, tt.test(t))
			defer m.Close()
			d := NewDiscordBot("", append([]DiscordOption{WithClient(m), WithBucket(mBucket)}, tt.extraOpts...)...)
			mal.Register(d)
			if tt.error != nil {
				_ = mal.Wrap(tt.error).Message(tt.message).Context(tt.context...).Notify(context.Background())
			} else {
				_ = mal.NewEntry(tt.message).Context(tt.context...).Notify(context.Background())
			}
			err := mal.Wait(context.Background())
			if err != nil {
				t.Fatalf("maleo.Wait() error = %v", err)
			}
			m.Wait()
			if m.count != tt.wantCount {
				t.Errorf("m.count = %d, want = %d", m.count, tt.wantCount)
			}
			tt.testBucket(t, mBucket)
		})
	}
}
