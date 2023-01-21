package maleohttp

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kinbiko/jsonassert"

	"github.com/tigorlazuardi/maleo"
)

var service = maleo.Service{
	Name:        "TestNewRoundTrip",
	Environment: "testing",
	Type:        "testing",
	Version:     "v0.1.0",
}

func TestNewRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		maleo   *maleo.Maleo
		handler http.Handler
		request func(server *httptest.Server) *http.Request
		test    func(t *testing.T, resp *http.Response, lg *maleo.TestingJSONLogger)
	}{
		{
			name:  "success",
			maleo: maleo.New(service),
			handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte("hello world"))
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}),
			request: func(server *httptest.Server) *http.Request {
				req, _ := http.NewRequest(http.MethodGet, server.URL, bytes.NewBufferString("hello"))
				return req
			},
			test: func(t *testing.T, _ *http.Response, lg *maleo.TestingJSONLogger) {
				j := jsonassert.New(t)
				want := `
				{
					"time": "<<PRESENCE>>",
					"code": 200,
					"message": "<<PRESENCE>>",
					"caller": "<<PRESENCE>>",
					"level": "info",
					"service": {
						"name": "TestNewRoundTrip",
						"environment": "testing",
						"type": "testing",
						"version": "v0.1.0"
					},
					"context": {
						"request": {
							"method": "GET",
							"url": "<<PRESENCE>>",
							"body": "hello"
						},
						"response": {
							"body": "hello world",
							"header": {
								"Content-Length": [
									"11"
								],
								"Content-Type": [
									"text/plain; charset=utf-8"
								],
								"Date": [
									"<<PRESENCE>>"
								]
							},
							"status": "200 OK"
						}
					}
				}`
				got := lg.String()
				j.Assertf(got, want)
				if !strings.Contains(got, "success: GET http://127.0.0.1:") {
					t.Errorf("expected message to contain success text, http method, and target url, got: %s", got)
				}
				if !strings.Contains(got, "round_trip_test.go:") {
					t.Error("expected to contain round_trip_test.go as caller")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := maleo.NewTestingJSONLogger()
			tt.maleo.SetLogger(logger)
			server := httptest.NewServer(tt.handler)
			defer server.Close()
			req := tt.request(server)
			rt := NewRoundTrip(Option.RoundTrip().Maleo(tt.maleo))
			client := &http.Client{Transport: rt}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if len(b) == 0 {
				t.Error("expected response body to be non-empty")
			}
			err = resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			tt.test(t, resp, logger)
			if t.Failed() {
				logger.PrettyPrint()
			}
		})
	}
}

func TestWrapHTTPClient_Get(t *testing.T) {
	type args struct {
		client *http.Client
		opts   []RoundTripOption
	}
	tests := []struct {
		name    string
		args    args
		maleo   *maleo.Maleo
		handler http.Handler
		test    func(t *testing.T, lg *maleo.TestingJSONLogger)
	}{
		{
			name:  "200 status code",
			args:  args{},
			maleo: maleo.New(service),
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusOK)
				_, err := writer.Write([]byte(`{"hello":"world"}`))
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}),
			test: func(t *testing.T, lg *maleo.TestingJSONLogger) {
				j := jsonassert.New(t)
				want := `
				{
					"time": "<<PRESENCE>>",
					"code": 200,
					"message": "<<PRESENCE>>",
					"caller": "<<PRESENCE>>",
					"level": "info",
					"service": {
						"name": "TestNewRoundTrip",
						"environment": "testing",
						"type": "testing",
						"version": "v0.1.0"
					},
					"context": {
						"request": {
							"method": "GET",
							"url": "<<PRESENCE>>"
						},
						"response": {
							"body": {"hello": "world"},
							"header": {
								"Content-Length": [
									"17"
								],
								"Content-Type": [
									"application/json"
								],
								"Date": [
									"<<PRESENCE>>"
								]
							},
							"status": "200 OK"
						}
					}
				}`
				got := lg.String()
				j.Assertf(got, want)
				if !strings.Contains(got, "success: GET http://127.0.0.1:") {
					t.Errorf("expected message to contain success text, http method, and target url, got: %s", got)
				}
				if !strings.Contains(got, "round_trip_test.go:") {
					t.Error("expected to contain this file as caller")
				}
			},
		},
		{
			name: "400 status code",
			args: args{
				opts: Option.RoundTrip().CallerDepth(5).AddCallerDepth(1),
			},
			maleo: maleo.New(service),
			handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusBadRequest)
				_, err := writer.Write([]byte(`{"hello":"world"}`))
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}),
			test: func(t *testing.T, lg *maleo.TestingJSONLogger) {
				j := jsonassert.New(t)
				want := `
				{
					"time": "<<PRESENCE>>",
					"code": 400,
					"message": "<<PRESENCE>>",
					"caller": "<<PRESENCE>>",
					"level": "error",
					"service": {
						"name": "TestNewRoundTrip",
						"environment": "testing",
						"type": "testing",
						"version": "v0.1.0"
					},
					"context": {
						"request": {
							"method": "GET",
							"url": "<<PRESENCE>>"
						},
						"response": {
							"body": {"hello": "world"},
							"header": {
								"Content-Length": [
									"17"
								],
								"Content-Type": [
									"application/json"
								],
								"Date": [
									"<<PRESENCE>>"
								]
							},
							"status": "400 Bad Request"
						}
					},
					"error": {
						"summary": "<<PRESENCE>>"
					}
				}`
				got := lg.String()
				j.Assertf(got, want)
				if !strings.Contains(got, "error: GET http://127.0.0.1:") {
					t.Error("expected message to contain success text, http method, and target url")
				}
				if !strings.Contains(got, "round_trip_test.go:") {
					t.Error("expected to contain this file as caller")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := maleo.NewTestingJSONLogger()
			tt.maleo.SetLogger(logger)
			server := httptest.NewServer(tt.handler)
			defer server.Close()
			opts := append(tt.args.opts, Option.RoundTrip().Maleo(tt.maleo))
			client := WrapHTTPClient(tt.args.client, opts...)
			resp, err := client.Get(server.URL)
			if err != nil {
				t.Fatal(err)
			}
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if len(b) == 0 {
				t.Error("expected response body to be non-empty")
			}
			err = resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			tt.test(t, logger)
			if t.Failed() {
				logger.PrettyPrint()
			}
		})
	}
}
