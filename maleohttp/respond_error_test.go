package maleohttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kinbiko/jsonassert"

	"github.com/tigorlazuardi/maleo"
)

type mockNullErrorTransformer struct{}

func (m mockNullErrorTransformer) ErrorBodyTransform(ctx context.Context, err error) any {
	return nil
}

type mockErrorEncoder struct{}

func (m mockErrorEncoder) ContentType() string {
	return ""
}

func (m mockErrorEncoder) Encode(input any) ([]byte, error) {
	return nil, errors.New("mock error encoder")
}

func TestResponder_RespondError(t *testing.T) {
	type fields struct {
		encoder          Encoder
		transformer      BodyTransformer
		errorTransformer ErrorBodyTransformer
		compressor       Compressor
		callerDepth      int
	}
	type testRequestGenerator = func(server *httptest.Server) *http.Request
	maleoGen := func(logger maleo.Logger) *maleo.Maleo {
		t := maleo.New(maleo.Service{
			Name:        "responder-test",
			Environment: "testing",
			Type:        "unit-test",
		})
		t.SetLogger(logger)
		return t
	}
	getRequest := func(server *httptest.Server) *http.Request {
		req, err := http.NewRequest(http.MethodGet, server.URL, nil)
		if err != nil {
			t.Fatal(err)
		}
		return req
	}
	postRequest := func(body io.ReadCloser) testRequestGenerator {
		return func(server *httptest.Server) *http.Request {
			req, err := http.NewRequest(http.MethodPost, server.URL, body)
			if err != nil {
				t.Fatal(err)
			}
			return req
		}
	}
	mustJsonBody := func(body any) io.ReadCloser {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		return io.NopCloser(bytes.NewReader(b))
	}
	tests := []struct {
		name    string
		fields  fields
		server  func(*Responder) *httptest.Server
		request func(server *httptest.Server) *http.Request
		test    func(t *testing.T, resp *http.Response, logger *maleo.TestingJSONLogger)
	}{
		{
			name: "common pattern",
			fields: fields{
				encoder:          NewJSONEncoder(),
				transformer:      NoopBodyTransform{},
				errorTransformer: SimpleErrorTransformer{},
				compressor:       NoCompression{},
				callerDepth:      2,
			},
			server: func(responder *Responder) *httptest.Server {
				handler := responder.RequestBodyCloner()(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					_, err := io.ReadAll(request.Body)
					if err != nil {
						t.Fatalf("failed to read request body: %v", err)
					}
					responder.RespondError(writer, request, errors.New("test error"))
				}))
				return httptest.NewServer(handler)
			},
			request: postRequest(mustJsonBody(map[string]any{"foo": "bar"})),
			test: func(t *testing.T, resp *http.Response, logger *maleo.TestingJSONLogger) {
				if resp.StatusCode != http.StatusInternalServerError {
					t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, resp.StatusCode)
				}
				if resp.Header.Get("Content-Type") != "application/json" {
					t.Errorf("expected content type %s, got %s", "application/json", resp.Header.Get("Content-Type"))
				}
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("failed to read response body: %v", err)
				}
				if len(body) == 0 {
					t.Error("expected response body, got empty")
				}
				wantBody := `{"error":"test error"}`
				j := jsonassert.New(t)
				j.Assertf(string(body), wantBody)
				wantLog := `
				{
					"time": "<<PRESENCE>>",
					"code": 500,
					"message": "test error",
					"caller": "<<PRESENCE>>",
					"level": "error",
					"service": {
						"name": "responder-test",
						"environment": "testing",
						"type": "unit-test"
					},
					"context": {
						"request": {
							"headers": {
								"Accept-Encoding": [
									"gzip"
								],
								"User-Agent": [
									"Go-http-client/1.1"
								]
							},
							"method": "POST",
							"url": "%s/",
							"body": {"foo":"bar"}
						},
						"response": {
							"body": {
								"error": "test error"
							},
							"headers": {
								"Content-Length": [
									"23"
								],
								"Content-Type": [
									"application/json"
								]
							},
							"status": 500
						}
					},
					"error": {
						"summary": "test error"
					}
				}`
				j.Assertf(logger.String(), wantLog, resp.Request.Host)
				if !strings.Contains(logger.String(), "maleohttp/respond_error_test.go") {
					t.Error("expected caller to be in maleohttp/respond_error_test.go")
				}
			},
		},
		{
			name: "compressed response",
			fields: fields{
				encoder:          NewJSONEncoder(),
				transformer:      NoopBodyTransform{},
				errorTransformer: SimpleErrorTransformer{},
				compressor:       NewGzipCompression(),
				callerDepth:      2,
			},
			server: func(responder *Responder) *httptest.Server {
				handler := responder.RequestBodyCloner()(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					_, err := io.ReadAll(request.Body)
					if err != nil {
						t.Fatalf("failed to read request body: %v", err)
					}
					responder.RespondError(writer, request, errors.New(strings.Repeat("test error ", 200)))
				}))
				return httptest.NewServer(handler)
			},
			request: postRequest(mustJsonBody(map[string]any{"foo": "bar"})),
			test: func(t *testing.T, resp *http.Response, logger *maleo.TestingJSONLogger) {
				if resp.StatusCode != http.StatusInternalServerError {
					t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, resp.StatusCode)
				}
				if resp.Header.Get("Content-Type") != "application/json" {
					t.Errorf("expected content type %s, got %s", "application/json", resp.Header.Get("Content-Type"))
				}
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("failed to read response body: %v", err)
				}
				if len(body) == 0 {
					t.Error("expected response body, got empty")
				}
				wantBody := `{"error":"<<PRESENCE>>"}`
				j := jsonassert.New(t)
				j.Assertf(string(body), wantBody)
				wantLog := `
				{
					"time": "<<PRESENCE>>",
					"code": 500,
					"message": "<<PRESENCE>>",
					"caller": "<<PRESENCE>>",
					"level": "error",
					"service": {
						"name": "responder-test",
						"environment": "testing",
						"type": "unit-test"
					},
					"context": {
						"request": {
							"headers": {
								"Accept-Encoding": [
									"gzip"
								],
								"User-Agent": [
									"Go-http-client/1.1"
								]
							},
							"method": "POST",
							"url": "%s/",
							"body": {"foo":"bar"}
						},
						"response": {
							"body": {
								"error": "<<PRESENCE>>"
							},
							"headers": {
								"Content-Encoding": ["gzip"],
								"Content-Length": [
									"62"
								],
								"Content-Type": [
									"application/json"
								]
							},
							"status": 500
						}
					},
					"error": {
						"summary": "<<PRESENCE>>"
					}
				}`
				j.Assertf(logger.String(), wantLog, resp.Request.Host)
				if !strings.Contains(logger.String(), "maleohttp/respond_error_test.go") {
					t.Error("expected caller to be in maleohttp/respond_error_test.go")
				}
			},
		},
		{
			name: "maleo error pattern",
			fields: fields{
				encoder:          NewJSONEncoder(),
				transformer:      NoopBodyTransform{},
				errorTransformer: SimpleErrorTransformer{},
				compressor:       NoCompression{},
				callerDepth:      2,
			},
			server: func(responder *Responder) *httptest.Server {
				handler := responder.RequestBodyCloner()(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					_, err := io.ReadAll(request.Body)
					if err != nil {
						t.Fatalf("failed to read request body: %v", err)
					}
					err = responder.maleo.Bail("test bail error").Code(http.StatusTeapot).Freeze()
					responder.RespondError(writer, request, err)
				}))
				return httptest.NewServer(handler)
			},
			request: postRequest(mustJsonBody(map[string]any{"foo": "bar"})),
			test: func(t *testing.T, resp *http.Response, logger *maleo.TestingJSONLogger) {
				if resp.StatusCode != http.StatusTeapot {
					t.Errorf("expected status code %d, got %d", http.StatusTeapot, resp.StatusCode)
				}
				if resp.Header.Get("Content-Type") != "application/json" {
					t.Errorf("expected content type %s, got %s", "application/json", resp.Header.Get("Content-Type"))
				}
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("failed to read response body: %v", err)
				}
				if len(body) == 0 {
					t.Error("expected response body, got empty")
				}
				wantBody := `{"error":"test bail error"}`
				j := jsonassert.New(t)
				j.Assertf(string(body), wantBody)
				wantLog := `
					{
						"time": "<<PRESENCE>>",
						"code": 418,
						"message": "test bail error",
						"caller": "<<PRESENCE>>",
						"level": "error",
						"service": {
							"name": "responder-test",
							"environment": "testing",
							"type": "unit-test"
						},
						"context": {
							"request": {
								"body": {
									"foo": "bar"
								},
								"headers": {
									"Accept-Encoding": [
										"gzip"
									],
									"User-Agent": [
										"Go-http-client/1.1"
									]
								},
								"method": "POST",
								"url": "%s/"
							},
							"response": {
								"body": {
									"error": "test bail error"
								},
								"headers": {
									"Content-Length": [
										"28"
									],
									"Content-Type": [
										"application/json"
									]
								},
								"status": 418
							}
						},
						"error": {
							"message": "test bail error",
							"caller": "<<PRESENCE>>",
							"error": {
								"summary": "test bail error"
							}
						}
					}`
				j.Assertf(logger.String(), wantLog, resp.Request.Host)
				if !strings.Contains(logger.String(), "maleohttp/respond_error_test.go") {
					t.Error("expected caller to be in maleohttp/respond_error_test.go")
				}
			},
		},
		{
			name: "handled nil error",
			fields: fields{
				encoder:          NewJSONEncoder(),
				transformer:      NoopBodyTransform{},
				errorTransformer: SimpleErrorTransformer{},
				compressor:       NoCompression{},
				callerDepth:      2,
			},
			server: func(responder *Responder) *httptest.Server {
				handler := responder.RequestBodyCloner()(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					_, err := io.ReadAll(request.Body)
					if err != nil {
						t.Fatalf("failed to read request body: %v", err)
					}
					responder.RespondError(writer, request, nil)
				}))
				return httptest.NewServer(handler)
			},
			request: postRequest(mustJsonBody(map[string]any{"foo": "bar"})),
			test: func(t *testing.T, resp *http.Response, logger *maleo.TestingJSONLogger) {
				if resp.StatusCode != http.StatusInternalServerError {
					t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, resp.StatusCode)
				}
				if resp.Header.Get("Content-Type") != "application/json" {
					t.Errorf("expected content type %s, got %s", "application/json", resp.Header.Get("Content-Type"))
				}
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("failed to read response body: %v", err)
				}
				if len(body) == 0 {
					t.Error("expected response body, got empty")
				}
				wantBody := `{"error":"Internal Server Error"}`
				j := jsonassert.New(t)
				j.Assertf(string(body), wantBody)
				wantLog := `
				{
					"time": "<<PRESENCE>>",
					"code": 500,
					"message": "Internal Server Error",
					"caller": "<<PRESENCE>>",
					"level": "error",
					"service": {
						"name": "responder-test",
						"environment": "testing",
						"type": "unit-test"
					},
					"context": {
						"request": {
							"headers": {
								"Accept-Encoding": [
									"gzip"
								],
								"User-Agent": [
									"Go-http-client/1.1"
								]
							},
							"method": "POST",
							"url": "%s/",
							"body": {"foo":"bar"}
						},
						"response": {
							"body": %s,
							"headers": {
								"Content-Length": [
									"<<PRESENCE>>"
								],
								"Content-Type": [
									"application/json"
								]
							},
							"status": 500
						}
					},
					"error": {
						"summary": "Internal Server Error",
						"details": "Internal Server Error"
					}
				}`
				j.Assertf(logger.String(), wantLog, resp.Request.Host, wantBody)
				if !strings.Contains(logger.String(), "maleohttp/respond_error_test.go") {
					t.Error("expected caller to be in maleohttp/respond_error_test.go")
				}
			},
		},
		{
			name: "still received body even when request body is another type of reader",
			fields: fields{
				encoder:          NewJSONEncoder(),
				transformer:      NoopBodyTransform{},
				errorTransformer: SimpleErrorTransformer{},
				compressor:       NoCompression{},
				callerDepth:      2,
			},
			server: func(responder *Responder) *httptest.Server {
				handler := responder.RequestBodyCloner()(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					request.Body = io.NopCloser(request.Body)
					_, err := io.ReadAll(request.Body)
					if err != nil {
						t.Fatalf("failed to read request body: %v", err)
					}
					responder.RespondError(writer, request, nil)
				}))
				return httptest.NewServer(handler)
			},
			request: postRequest(mustJsonBody(map[string]any{"foo": "bar"})),
			test: func(t *testing.T, resp *http.Response, logger *maleo.TestingJSONLogger) {
				if resp.StatusCode != http.StatusInternalServerError {
					t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, resp.StatusCode)
				}
				if resp.Header.Get("Content-Type") != "application/json" {
					t.Errorf("expected content type %s, got %s", "application/json", resp.Header.Get("Content-Type"))
				}
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("failed to read response body: %v", err)
				}
				if len(body) == 0 {
					t.Error("expected response body, got empty")
				}
				wantBody := `{"error":"Internal Server Error"}`
				j := jsonassert.New(t)
				j.Assertf(string(body), wantBody)
				wantLog := `
				{
					"time": "<<PRESENCE>>",
					"code": 500,
					"message": "Internal Server Error",
					"caller": "<<PRESENCE>>",
					"level": "error",
					"service": {
						"name": "responder-test",
						"environment": "testing",
						"type": "unit-test"
					},
					"context": {
						"request": {
							"headers": {
								"Accept-Encoding": [
									"gzip"
								],
								"User-Agent": [
									"Go-http-client/1.1"
								]
							},
							"method": "POST",
							"url": "%s/",
							"body": {"foo":"bar"}
						},
						"response": {
							"body": %s,
							"headers": {
								"Content-Length": [
									"<<PRESENCE>>"
								],
								"Content-Type": [
									"application/json"
								]
							},
							"status": 500
						}
					},
					"error": {
						"summary": "Internal Server Error",
						"details": "Internal Server Error"
					}
				}`
				j.Assertf(logger.String(), wantLog, resp.Request.Host, wantBody)
				if !strings.Contains(logger.String(), "maleohttp/respond_error_test.go") {
					t.Error("expected caller to be in maleohttp/respond_error_test.go")
				}
			},
		},
		{
			name: "no body is sent when transformer returns nil",
			fields: fields{
				encoder:          NewJSONEncoder(),
				transformer:      NoopBodyTransform{},
				errorTransformer: mockNullErrorTransformer{},
				compressor:       NoCompression{},
				callerDepth:      2,
			},
			server: func(responder *Responder) *httptest.Server {
				handler := responder.RequestBodyCloner()(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					request.Body = io.NopCloser(request.Body)
					_, err := io.ReadAll(request.Body)
					if err != nil {
						t.Fatalf("failed to read request body: %v", err)
					}
					responder.RespondError(writer, request, errors.New("foo"))
				}))
				return httptest.NewServer(handler)
			},
			request: postRequest(mustJsonBody(map[string]any{"foo": "bar"})),
			test: func(t *testing.T, resp *http.Response, logger *maleo.TestingJSONLogger) {
				if resp.StatusCode != http.StatusInternalServerError {
					t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, resp.StatusCode)
				}
				if resp.Header.Get("Content-Type") != "" {
					t.Errorf("expected no content type")
				}
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("failed to read response body: %v", err)
				}
				if len(body) != 0 {
					t.Errorf("expected no response body, got %s", string(body))
				}
				j := jsonassert.New(t)
				wantLog := `
				{
					"time": "<<PRESENCE>>",
					"code": 500,
					"message": "foo",
					"caller": "<<PRESENCE>>",
					"level": "error",
					"service": {
						"name": "responder-test",
						"environment": "testing",
						"type": "unit-test"
					},
					"context": {
						"request": {
							"headers": {
								"Accept-Encoding": [
									"gzip"
								],
								"User-Agent": [
									"Go-http-client/1.1"
								]
							},
							"method": "POST",
							"url": "%s/",
							"body": {"foo":"bar"}
						},
						"response": {
							"status": 500
						}
					},
					"error": {
						"summary": "foo"
					}
				}`
				j.Assertf(logger.String(), wantLog, resp.Request.Host)
				if !strings.Contains(logger.String(), "maleohttp/respond_error_test.go") {
					t.Error("expected caller to be in maleohttp/respond_error_test.go")
				}
			},
		},
		{
			name: "expected output on error encoding",
			fields: fields{
				encoder:          mockErrorEncoder{},
				transformer:      NoopBodyTransform{},
				errorTransformer: SimpleErrorTransformer{},
				compressor:       NoCompression{},
				callerDepth:      2,
			},
			server: func(responder *Responder) *httptest.Server {
				handler := responder.RequestBodyCloner()(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					request.Body = io.NopCloser(request.Body)
					_, err := io.ReadAll(request.Body)
					if err != nil {
						t.Fatalf("failed to read request body: %v", err)
					}
					responder.RespondError(writer, request, errors.New("foo"))
				}))
				return httptest.NewServer(handler)
			},
			request: getRequest,
			test: func(t *testing.T, resp *http.Response, logger *maleo.TestingJSONLogger) {
				if resp.StatusCode != http.StatusInternalServerError {
					t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, resp.StatusCode)
				}
				if resp.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
					t.Errorf("expected content type text/plain; charset=utf-8, got %s", resp.Header.Get("Content-Type"))
				}
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("failed to read response body: %v", err)
				}
				if len(body) == 0 {
					t.Errorf("expected response body")
				}
				j := jsonassert.New(t)
				wantLog := `
				{
					"time": "<<PRESENCE>>",
					"code": 500,
					"message": "foo",
					"caller": "<<PRESENCE>>",
					"level": "error",
					"service": {
						"name": "responder-test",
						"environment": "testing",
						"type": "unit-test"
					},
					"context": {
						"request": {
							"headers": {
								"Accept-Encoding": [
									"gzip"
								],
								"User-Agent": [
									"Go-http-client/1.1"
								]
							},
							"method": "GET",
							"url": "%s/"
						},
						"response": {
							"status": 500,
							"body": "ENCODING ERROR",
							"headers": {
								"Content-Type": [
									"text/plain; charset=utf-8"
								]
							}
						}
					},
					"error": {
						"summary": "foo"
					}
				}`
				j.Assertf(logger.String(), wantLog, resp.Request.Host)
				if !strings.Contains(logger.String(), "maleohttp/respond_error_test.go") {
					t.Error("expected caller to be in maleohttp/respond_error_test.go")
				}
			},
		},
		{
			name: "still received proper response even compressor fails",
			fields: fields{
				encoder: func() Encoder {
					enc := NewJSONEncoder()
					enc.SetIndent("")
					enc.SetPrefix("")
					enc.SetHtmlEscape(false)
					return enc
				}(),
				transformer:      NoopBodyTransform{},
				errorTransformer: SimpleErrorTransformer{},
				compressor:       mockErrorCompressor{},
				callerDepth:      2,
			},
			server: func(responder *Responder) *httptest.Server {
				handler := responder.RequestBodyCloner()(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					request.Body = io.NopCloser(request.Body)
					_, err := io.ReadAll(request.Body)
					if err != nil {
						t.Fatalf("failed to read request body: %v", err)
					}
					responder.RespondError(writer, request, nil)
				}))
				return httptest.NewServer(handler)
			},
			request: postRequest(mustJsonBody(map[string]any{"foo": "bar"})),
			test: func(t *testing.T, resp *http.Response, logger *maleo.TestingJSONLogger) {
				if resp.StatusCode != http.StatusInternalServerError {
					t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, resp.StatusCode)
				}
				if resp.Header.Get("Content-Type") != "application/json" {
					t.Errorf("expected content type %s, got %s", "application/json", resp.Header.Get("Content-Type"))
				}
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("failed to read response body: %v", err)
				}
				if len(body) == 0 {
					t.Error("expected response body, got empty")
				}
				wantBody := `{"error":"Internal Server Error"}`
				j := jsonassert.New(t)
				j.Assertf(string(body), wantBody)
				logs := strings.Split(logger.String(), "\n")
				// 3 split because EOL is \n
				if len(logs) != 3 {
					t.Fatalf("expected 2 logs, got %d", len(logs))
				}
				wantEntry := `
				{
					"time": "<<PRESENCE>>",
					"code": 500,
					"message": "compress error",
					"caller": "<<PRESENCE>>",
					"level": "warn",
					"service": {
						"name": "responder-test",
						"environment": "testing",
						"type": "unit-test"
					},
					"error": {
						"summary": "compress error"
					}
				}`
				j.Assertf(logs[0], wantEntry)
				wantLog := `
				{
					"time": "<<PRESENCE>>",
					"code": 500,
					"message": "Internal Server Error",
					"caller": "<<PRESENCE>>",
					"level": "error",
					"service": {
						"name": "responder-test",
						"environment": "testing",
						"type": "unit-test"
					},
					"context": {
						"request": {
							"headers": {
								"Accept-Encoding": [
									"gzip"
								],
								"User-Agent": [
									"Go-http-client/1.1"
								]
							},
							"method": "POST",
							"url": "%s/",
							"body": {"foo":"bar"}
						},
						"response": {
							"body": %s,
							"headers": {
								"Content-Length": [
									"<<PRESENCE>>"
								],
								"Content-Type": [
									"application/json"
								]
							},
							"status": 500
						}
					},
					"error": {
						"summary": "Internal Server Error",
						"details": "Internal Server Error"
					}
				}`
				j.Assertf(logs[1], wantLog, resp.Request.Host, wantBody)
				if !strings.Contains(logger.String(), "maleohttp/respond_error_test.go") {
					t.Error("expected caller to be in maleohttp/respond_error_test.go")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := maleo.NewTestingJSONLogger()
			mal := maleoGen(logger)
			r := NewResponder()
			r.SetEncoder(tt.fields.encoder)
			r.SetBodyTransformer(tt.fields.transformer)
			r.SetErrorTransformer(tt.fields.errorTransformer)
			r.SetMaleo(mal)
			r.SetCompressor(tt.fields.compressor)
			r.SetCallerDepth(tt.fields.callerDepth)
			r.RegisterHook(NewLoggerHook())
			server := tt.server(r)
			defer server.Close()
			resp, err := http.DefaultClient.Do(tt.request(server))
			if err != nil {
				t.Fatal(err)
			}
			tt.test(t, resp, logger)
			err = resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			if t.Failed() {
				logger.PrettyPrint()
			}
		})
	}
}
