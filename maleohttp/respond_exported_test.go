package maleohttp_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/tigorlazuardi/maleo"
	"github.com/tigorlazuardi/maleo/maleohttp"
)

func TestGlobalRespond(t *testing.T) {
	const envKey = "TOWER_HTTP_TEST_EXPORTED"
	if os.Getenv(envKey) == "" {
		t.Skipf("skipping test; set %s env to run", envKey)
	}
	maleoGen := func(logger maleo.Logger) *maleo.Maleo {
		t := maleo.NewMaleo(maleo.Service{
			Name:        "responder-test",
			Environment: "testing",
			Type:        "unit-test",
		})
		t.SetLogger(logger)
		return t
	}
	tests := []struct {
		name   string
		server func() *httptest.Server
		test   func(t *testing.T, logger *maleo.TestingJSONLogger)
	}{
		{
			name: "expected caller location for respond",
			server: func() *httptest.Server {
				handler := maleohttp.RequestBodyCloner()(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					maleohttp.Respond(writer, request, nil)
				}))
				return httptest.NewServer(handler)
			},
			test: func(t *testing.T, logger *maleo.TestingJSONLogger) {
				if !strings.Contains(logger.String(), "maleohttp/respond_exported_test.go") {
					t.Error("expected caller location is correct")
				}
			},
		},
		{
			name: "expected caller location for respond error",
			server: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					maleohttp.RespondError(w, r, nil)
				}))
			},
			test: func(t *testing.T, logger *maleo.TestingJSONLogger) {
				if !strings.Contains(logger.String(), "maleohttp/respond_exported_test.go") {
					t.Error("expected caller location is correct")
				}
			},
		},
		{
			name: "expected caller location for respond stream",
			server: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					maleohttp.RespondStream(w, r, "", nil)
				}))
			},
			test: func(t *testing.T, logger *maleo.TestingJSONLogger) {
				if !strings.Contains(logger.String(), "maleohttp/respond_exported_test.go") {
					t.Error("expected caller location is correct")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := maleo.NewTestingJSONLogger()
			tow := maleoGen(logger)
			r := maleohttp.NewResponder()
			r.SetMaleo(tow)
			r.RegisterHook(maleohttp.NewLoggerHook())
			r.SetCallerDepth(3)
			maleohttp.Exported.Responder().SetMaleo(tow)
			maleohttp.Exported.Responder().RegisterHook(maleohttp.NewLoggerHook())
			maleohttp.Exported.SetResponder(r)
			server := tt.server()
			defer server.Close()
			req, err := http.NewRequest(http.MethodGet, server.URL, nil)
			if err != nil {
				t.Errorf("Error creating request: %s", err.Error())
				return
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			tt.test(t, logger)
			err = resp.Body.Close()
			if err != nil {
				t.Fatalf("Error closing response body: %s", err.Error())
			}
			if t.Failed() {
				logger.PrettyPrint()
			}
		})
	}
}
