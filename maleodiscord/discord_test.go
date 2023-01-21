package maleodiscord

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/kinbiko/jsonassert"

	"github.com/tigorlazuardi/maleo"
	"github.com/tigorlazuardi/maleo/bucket"
	"github.com/tigorlazuardi/maleo/locker"
)

type testHook struct {
	t                  *testing.T
	wg                 *sync.WaitGroup
	checkBucketContext bool
}

func (t testHook) PreMessageHook(ctx context.Context, _ *WebhookContext) context.Context {
	ctx = context.WithValue(ctx, "test", "test")
	return ctx
}

func (t testHook) PostMessageHook(ctx context.Context, _ *WebhookContext, err error) {
	defer t.wg.Done()
	if err != nil {
		t.t.Error(err)
	}
	if e, ok := ctx.Value("test").(string); ok {
		if e != "test" {
			t.t.Errorf("context value of test should have value of 'test' not '%s'", e)
		}
	} else {
		t.t.Error("context value of test should exist in PostMessageHook")
	}
	if t.checkBucketContext {
		if e, ok := ctx.Value("test-bucket").(string); ok {
			if e != "test-bucket" {
				t.t.Errorf("context value of test should have value of 'test-bucket' not '%s'", e)
			}
		} else {
			t.t.Error("context value of test-bucket should exist in PostMessageHook")
		}
	}
}

func (t testHook) PreBucketUploadHook(ctx context.Context, _ *WebhookContext) context.Context {
	if e, ok := ctx.Value("test").(string); ok {
		if e != "test" {
			t.t.Errorf("context value of test should have value of 'test' not '%s'", e)
		}
	} else {
		t.t.Error("context value of test should exist in PreBucketUploadHook")
	}
	ctx = context.WithValue(ctx, "test-bucket", "test-bucket")
	return ctx
}

func (t testHook) PostBucketUploadHook(ctx context.Context, _ *WebhookContext, results []bucket.UploadResult) {
	defer t.wg.Done()
	for _, result := range results {
		if result.Error != nil {
			t.t.Error(result.Error)
		}
	}
	if e, ok := ctx.Value("test").(string); ok {
		if e != "test" {
			t.t.Errorf("context value of test should have value of 'test' not '%s'", e)
		}
	} else {
		t.t.Error("context value of test should exist in PostBucketUploadHook")
	}
	if e, ok := ctx.Value("test-bucket").(string); ok {
		if e != "test-bucket" {
			t.t.Errorf("context value of test should have value of 'test-bucket' not '%s'", e)
		}
	} else {
		t.t.Error("context value of test-bucket should exist in PostBucketUploadHook")
	}
}

type mockClient struct {
	*http.Client
	*httptest.Server
	count int
	*sync.WaitGroup
	test *testing.T
}

type callback = func(r *http.Request)

func (m mockClient) Do(req *http.Request) (*http.Response, error) {
	var err error
	req.URL, err = url.Parse(m.Server.URL)
	if err != nil {
		m.test.Errorf("failed to parse url: %v", err)
		return nil, err
	}
	return m.Client.Do(req)
}

func newMockClient(t *testing.T, cb callback) *mockClient {
	m := &mockClient{test: t}
	m.WaitGroup = &sync.WaitGroup{}
	m.WaitGroup.Add(1)
	m.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cb != nil {
			cb(r)
		}
		w.WriteHeader(http.StatusNoContent)
		m.count += 1
		m.WaitGroup.Done()
	}))
	m.Client = http.DefaultClient
	return m
}

func formatJSON(b []byte) string {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "\t")
	if err != nil {
		return string(b)
	}
	return out.String()
}

func TestDiscord(t *testing.T) {
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
			name: "expected - entry",
			test: func(t *testing.T) callback {
				return func(r *http.Request) {
					if r.Method != http.MethodPost {
						t.Errorf("expected method to be %s, got %s", http.MethodPost, r.Method)
					}
					if r.Header.Get("Content-Type") != "application/json" {
						t.Errorf("expected content type to be application/json, got %s", r.Header.Get("Content-Type"))
					}
					body, err := io.ReadAll(r.Body)
					if err != nil {
						t.Fatalf("failed to read body: %v", err)
					}
					j := jsonassert.New(t)
					want := `
					{
						"content": "@here Message from service **test** on type **test** on environment **test**",
						"embeds": [
							{
								"title": "Summary",
								"type": "rich",
								"description": "<<PRESENCE>>",
								"color": 1606980
							},
							{
								"title": "Context",
								"type": "rich",
								"description": "<<PRESENCE>>",
								"color": 407920
							},
							{
								"title": "Metadata",
								"type": "rich",
								"description": "<<PRESENCE>>",
								"timestamp": "<<PRESENCE>>",
								"color": 6576731,
								"fields": [
									{
										"name": "foo",
										"value": "bar",
										"inline": true
									},
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
							}
						]
					}`
					j.Assertf(string(body), want)
					if t.Failed() {
						fmt.Println(formatJSON(body))
					}
				}
			},
			wantCount: 1,
			context:   []any{maleo.F{"foo": "bar"}},
			error:     nil,
			message:   "foo",
			extraOpts: []DiscordOption{
				WithName("test-discord"),
				WithLock(locker.NewLocalLock()),
				WithSemaphore(make(chan struct{}, 3)),
				WithTrace(maleo.TraceCapturerFunc(func(ctx context.Context) []maleo.KVString {
					return []maleo.KVString{{Key: "foo", Value: "bar"}}
				})),
				WithBucket(nil),
				WithGlobalKey("global"),
				WithCooldown(time.Second),
				WithDataEncoder(JSONDataEncoder{}),
				WithCodeBlockBuilder(JSONCodeBlockBuilder{}),
			},
		},
		{
			name: "expected - error",
			test: func(t *testing.T) callback {
				return func(r *http.Request) {
					if r.Method != http.MethodPost {
						t.Errorf("expected method to be %s, got %s", http.MethodPost, r.Method)
					}
					if r.Header.Get("Content-Type") != "application/json" {
						t.Errorf("expected content type to be application/json, got %s", r.Header.Get("Content-Type"))
					}
					body, err := io.ReadAll(r.Body)
					if err != nil {
						t.Fatalf("failed to read body: %v", err)
					}
					j := jsonassert.New(t)
					want := `
					{
						"content": "@here an error has occurred on service **test** on type **test** on environment **test**",
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
								"title": "Context",
								"type": "rich",
								"description": "<<PRESENCE>>",
								"color": 407920
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
							}
						]
					}`
					j.Assertf(string(body), want)
					if t.Failed() {
						fmt.Println(formatJSON(body))
					}
				}
			},
			wantCount: 1,
			context:   []any{maleo.F{"foo": "bar"}, "1000"},
			error:     errors.New("bar"),
			message:   "foo",
			extraOpts: []DiscordOption{WithHook(func() Hook {
				wg := &sync.WaitGroup{}
				wg.Add(1)
				return testHook{
					t:                  t,
					wg:                 wg,
					checkBucketContext: false,
				}
			}())},
		},
		{
			name: "with custom embed builder",
			test: func(t *testing.T) callback {
				return func(r *http.Request) {
					if r.Method != http.MethodPost {
						t.Errorf("expected method to be %s, got %s", http.MethodPost, r.Method)
					}
					if r.Header.Get("Content-Type") != "application/json" {
						t.Errorf("expected content type to be application/json, got %s", r.Header.Get("Content-Type"))
					}
					body, err := io.ReadAll(r.Body)
					if err != nil {
						t.Fatalf("failed to read body: %v", err)
					}
					j := jsonassert.New(t)
					want := `
					{
						"content": "@here an error has occurred on service **test** on type **test** on environment **test**"
					}`
					j.Assertf(string(body), want)
					if t.Failed() {
						fmt.Println(formatJSON(body))
					}
				}
			},
			wantCount: 1,
			context:   []any{maleo.F{"foo": "bar"}, "1000"},
			error:     errors.New("bar"),
			message:   "foo",
			extraOpts: []DiscordOption{WithEmbedBuilder(EmbedBuilderFunc(func(ctx context.Context, msg maleo.MessageContext, info *ExtraInformation) ([]*Embed, []bucket.File) {
				return []*Embed{}, []bucket.File{}
			}))},
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
