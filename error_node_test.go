package maleo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/kinbiko/jsonassert"
)

func TestErrorNode_CodeBlockJSON(t *testing.T) {
	tests := []struct {
		name      string
		baseError error
		messages  []string
		want      string
		wantErr   bool
	}{
		{
			name:      "expected output",
			baseError: errors.New("base error"),
			messages:  []string{"message 1", "message 2", "message 3"},
			want: `
{
   "message": "message 3",
   "caller": "<<PRESENCE>>",
   "error": {
      "message": "message 2",
      "caller": "<<PRESENCE>>",
      "error": {
         "message": "message 1",
         "caller": "<<PRESENCE>>",
         "error": {
            "time": "<<PRESENCE>>",
            "code": 500,
            "message": "base error",
            "caller": "<<PRESENCE>>",
            "level": "error",
            "service": {
               "name": "test",
               "environment": "test",
               "type": "test",
			   "version": "v0.1.0-test"
            },
            "error": {
               "summary": "base error"
            }
         }
      }
   }
}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mal, _ := NewTestingMaleo()
			err := mal.Wrap(tt.baseError).Freeze()
			for _, e := range tt.messages {
				err = mal.WrapFreeze(err, e)
			}
			got, errCB := err.(*ErrorNode).CodeBlockJSON()
			if (errCB != nil) != tt.wantErr {
				t.Errorf("ErrorNode.CodeBlockJSON() error = %v, wantErr %v", errCB, tt.wantErr)
				return
			}
			j := jsonassert.New(t)
			j.Assertf(string(got), tt.want)
			if t.Failed() {
				fmt.Println(string(got))
			}
			if !strings.Contains(string(got), "error_node_test.go") {
				t.Error("expected to see caller in error_node_test.go")
			}
			if strings.Count(string(got), "error_node_test.go") != 4 {
				t.Error("expected to see four callers field in error_node_test.go")
			}
		})
	}
}

type mockImplError struct{}

func (m mockImplError) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"message": m.Message(),
		"caller":  m.Caller(),
		"code":    m.Code(),
		"error":   m.Error(),
		"level":   m.Level(),
		"service": m.Service(),
		"time":    m.Time(),
	})
}

func (m mockImplError) Error() string {
	return "mock"
}

func (m mockImplError) Caller() Caller {
	return GetCaller(1)
}

func (m mockImplError) Code() int {
	return 500
}

func (m mockImplError) Context() []any {
	return nil
}

func (m mockImplError) Unwrap() error {
	return nil
}

func (m mockImplError) WriteError(w LineWriter) {}

func (m mockImplError) HTTPCode() int {
	return 500
}

func (m mockImplError) Key() string {
	return ""
}

func (m mockImplError) Level() Level {
	return ErrorLevel
}

func (m mockImplError) Message() string {
	return "mocking time"
}

func (m mockImplError) Time() time.Time {
	return time.Now()
}

func (m mockImplError) Service() Service {
	return Service{
		Name: "mock",
	}
}

func (m mockImplError) Log(ctx context.Context) Error {
	return m
}

func (m mockImplError) Notify(ctx context.Context, opts ...MessageOption) Error {
	return m
}

func TestErrorNode_MarshalJSON(t *testing.T) {
	mal, _ := NewTestingMaleo()
	tests := []struct {
		name    string
		err     *ErrorNode
		want    string
		wantErr bool
	}{
		{
			name: "expected output - simple error",
			err:  mal.Wrap(errors.New("base error")).Message("bar").Freeze().(*ErrorNode),
			want: `
				{
				   "time": "<<PRESENCE>>",
				   "code": 500,
				   "message": "bar",
				   "caller": "<<PRESENCE>>",
				   "level": "error",
				   "service": {
					  "name": "test",
					  "environment": "test",
					  "type": "test",
					  "version": "v0.1.0-test"
				   },
				   "error": {
					  "summary": "base error"
				   }
				}`,
			wantErr: false,
		},
		{
			name: "expected output - nested error",
			err: func() *ErrorNode {
				base := mal.Wrap(errors.New("base error")).Message("bar").Freeze()
				return mal.WrapFreeze(base, "error 1").(*ErrorNode)
			}(),
			want: `
				{
				   "time": "<<PRESENCE>>",
				   "code": 500,
				   "message": "error 1",
				   "caller": "<<PRESENCE>>",
				   "level": "error",
				   "service": "<<PRESENCE>>",
				   "error": {
					  "message": "bar",
					  "caller": "<<PRESENCE>>",
					  "error": {
						 "summary": "base error"
					  }
				   }
				}`,
			wantErr: false,
		},
		{
			name: "expected output - nested error with wrap that does nothing",
			err: func() *ErrorNode {
				base := mal.Wrap(errors.New("base error")).Message("bar").Freeze()
				base = mal.Wrap(base).Freeze()
				base = mal.Wrap(base).Freeze()
				return mal.WrapFreeze(base, "error 1").(*ErrorNode)
			}(),
			want: `
				{
				   "time": "<<PRESENCE>>",
				   "code": 500,
				   "message": "error 1",
				   "caller": "<<PRESENCE>>",
				   "level": "error",
				   "service": {
					  "name": "test",
					  "environment": "test",
					  "type": "test",
					  "version": "v0.1.0-test"
				   },
				   "error": {
					  "message": "bar",
					  "caller": "<<PRESENCE>>",
					  "error": {
						 "summary": "base error"
					  }
				   }
				}`,
			wantErr: false,
		},
		{
			name: "expected output - wrap other Error implementation",
			err: func() *ErrorNode {
				base := mal.Wrap(mockImplError{}).Code(400).Freeze()
				return mal.WrapFreeze(base, "error 1").(*ErrorNode)
			}(),
			want: `
				{
				   "time": "<<PRESENCE>>",
				   "code": 400,
				   "message": "error 1",
				   "caller": "<<PRESENCE>>",
				   "level": "error",
				   "service": "<<PRESENCE>>",
				   "error": {
					  "caller": "<<PRESENCE>>",
					  "code": 500,
					  "error": "mock",
					  "level": 2,
					  "message": "mocking time",
					  "service": "<<PRESENCE>>",
					  "time": "<<PRESENCE>>"
				   }
				}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.err
			got, err := e.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			j := jsonassert.New(t)
			j.Assertf(string(got), tt.want)
			if t.Failed() {
				out := new(bytes.Buffer)
				_ = json.Indent(out, got, "", "   ")
				fmt.Println(out.String())
			}
		})
	}
}

func Test_Error(t *testing.T) {
	mal, _ := NewTestingMaleo()
	l := newMockLogger()
	mal.SetLogger(l)
	m := newMockMessenger(1)
	mal.Register(m)
	base := errors.New("based")
	builder := mal.Wrap(base).Message("message 1")
	err := builder.Freeze()
	if err.Error() != "message 1: based" {
		t.Errorf("Error.Error() = %v, want %v", err.Error(), "message 1")
	}
	if err.Message() != "message 1" {
		t.Errorf("Error.Message() = %v, want %v", err.Message(), "message 1")
	}
	if err.Code() != 500 {
		t.Errorf("Error.Code() = %v, want %v", err.Code(), 500)
	}
	if err.HTTPCode() != 500 {
		t.Errorf("Error.HTTPCode() = %v, want %v", err.HTTPCode(), 500)
	}
	if errors.Unwrap(err) != base {
		t.Errorf("Error.Unwrap() = %v, want %v", errors.Unwrap(err), base)
	}
	ctx := context.Background()
	_ = err.Log(ctx).Notify(ctx)
	if !l.called {
		t.Error("Expected logger to be called")
	}
	err2 := m.Wait(ctx)
	if err2 != nil {
		t.Fatalf("Expected messenger to wait without error, got %v", err2)
	}
	if !m.called {
		t.Error("Expected messenger to be called")
	}

	b, errMarshal := json.Marshal(err)
	if errMarshal != nil {
		t.Fatalf("Expected error to marshal to JSON without error, got %v", errMarshal)
	}
	defer func() {
		if t.Failed() {
			out := new(bytes.Buffer)
			_ = json.Indent(out, b, "", "  ")
			t.Log(out.String())
		}
	}()
	j := jsonassert.New(t)
	j.Assertf(string(b), `
	{
		"time": "<<PRESENCE>>",
		"code": 500,
		"message": "message 1",
		"caller": "<<PRESENCE>>",
		"level": "error",
		"service": "<<PRESENCE>>",
		"error": {"summary": "based"}
	}`)

	now := time.Now()
	base2 := errors.New("based 2")
	builder.Level(WarnLevel).
		Caller(GetCaller(1)).
		Code(600).
		Error(base2).
		Context(1).
		Key("foo").
		Time(now)

	err = builder.Freeze()
	if err.HTTPCode() != 500 {
		t.Errorf("Error.HTTPCode() = %v, want %v", err.HTTPCode(), 500)
	}
	builder.Code(1400)
	err = builder.Freeze()
	if err.HTTPCode() != 400 {
		t.Errorf("Error.HTTPCode() = %v, want %v", err.HTTPCode(), 400)
	}
	b, errMarshal = json.Marshal(err)
	if errMarshal != nil {
		t.Fatalf("Expected error to marshal to JSON without error, got %v", errMarshal)
	}
	j.Assertf(string(b), `
	{
		"time": "%s",
		"code": 1400,
		"message": "message 1",
		"caller": "<<PRESENCE>>",
		"key": "foo",
		"level": "warn",
		"service": "<<PRESENCE>>",
		"context": 1,
		"error": {"summary": "based 2"}
	}`, now.Format(time.RFC3339))

	builder.Context(2).Message("foo %s", "bar").Key("foo %s", "bar")
	err = builder.Freeze()
	b, errMarshal = json.Marshal(err)
	if errMarshal != nil {
		t.Fatalf("Expected error to marshal to JSON without error, got %v", errMarshal)
	}
	j.Assertf(string(b), `
	{
		"time": "%s",
		"code": 1400,
		"message": "foo bar",
		"caller": "<<PRESENCE>>",
		"key": "foo bar",
		"level": "warn",
		"service": "<<PRESENCE>>",
		"context": [1, 2],
		"error": {"summary": "based 2"}
	}`, now.Format(time.RFC3339))

	l2 := newMockLogger()
	m2 := newMockMessenger(1)
	mal.defaultParams.Messengers = []Messenger{}
	mal.Register(m2)
	mal.SetLogger(l2)
	_ = builder.Log(ctx)
	if !l2.called {
		t.Error("Expected logger to be called")
	}
	_ = builder.Notify(ctx)
	err3 := m2.Wait(ctx)
	if err3 != nil {
		t.Fatalf("Expected messenger to wait without error, got %v", err3)
	}
	if !m2.called {
		t.Error("Expected messenger to be called")
	}
}

type mockError struct {
	Message string
}

func (m mockError) Error() string {
	return m.Message
}

type mockJSONError struct{}

func (m mockJSONError) Error() string {
	return "mock"
}

func (m mockJSONError) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"foo": "bar"})
}

type funcError func()

func (f funcError) Error() string {
	return "mock"
}

func TestError_WrapMarshalJSON(t *testing.T) {
	mal, _ := NewTestingMaleo()
	base := mockError{Message: "based"}
	builder := mal.Wrap(base).Message("message 1")
	err := builder.Freeze()
	b, errMarshal := json.Marshal(err)
	if errMarshal != nil {
		t.Fatalf("Expected error to marshal to JSON without error, got %v", errMarshal)
	}
	defer func() {
		if t.Failed() {
			out := new(bytes.Buffer)
			_ = json.Indent(out, b, "", "    ")
			t.Log(out.String())
		}
	}()
	j := jsonassert.New(t)
	j.Assertf(string(b), `
	{
		"time": "<<PRESENCE>>",
		"code": 500,
		"message": "message 1",
		"caller": "<<PRESENCE>>",
		"level": "error",
		"service": "<<PRESENCE>>",
		"error": {"summary": "based", "details": {"Message": "based"}}
	}`)
	builder.Error(mockJSONError{})
	err = builder.Freeze()
	b, errMarshal = json.Marshal(err)
	if errMarshal != nil {
		t.Fatalf("Expected error to marshal to JSON without error, got %v", errMarshal)
	}
	j.Assertf(string(b), `
	{
		"time": "<<PRESENCE>>",
		"code": 500,
		"message": "message 1",
		"caller": "<<PRESENCE>>",
		"level": "error",
		"service": "<<PRESENCE>>",
		"error": {"foo": "bar"}
	}`)
	builder.Error(funcError(func() {}))
	err = builder.Freeze()
	b, errMarshal = json.Marshal(err)
	if errMarshal != nil {
		t.Fatalf("Expected error to marshal to JSON without error, got %v", errMarshal)
	}
	j.Assertf(string(b), `
	{
		"time": "<<PRESENCE>>",
		"code": 500,
		"message": "message 1",
		"caller": "<<PRESENCE>>",
		"level": "error",
		"service": "<<PRESENCE>>",
		"error": "mock"
	}`)
}

func Test_Error_WriteError(t *testing.T) {
	tests := []struct {
		name   string
		error  Error
		writer func() (LineWriter, fmt.Stringer)
		want   string
	}{
		{
			name: "No Duplicates",
			error: func() Error {
				err := BailFreeze("bail")
				err = WrapFreeze(err, "wrap")
				return Wrap(err).Freeze()
			}(),
			writer: func() (LineWriter, fmt.Stringer) {
				s := &strings.Builder{}
				lw := NewLineWriter(s).LineBreak(" => ").Build()
				return lw, s
			},
			want: "wrap => bail",
		},
		{
			name: "No Duplicates - Tail",
			error: func() Error {
				err := errors.New("errors.New")
				err = WrapFreeze(err, "wrap")
				err = Wrap(err).Freeze()
				return Wrap(err).Message("foo").Freeze()
			}(),
			writer: func() (LineWriter, fmt.Stringer) {
				s := &strings.Builder{}
				lw := NewLineWriter(s).LineBreak(" => ").Build()
				return lw, s
			},
			want: "foo => wrap => errors.New",
		},
		{
			name: "Ensure different messages are written",
			error: func() Error {
				err := BailFreeze("bail")
				err = WrapFreeze(err, "wrap")
				return Wrap(err).Message("wrap 2").Freeze()
			}(),
			writer: func() (LineWriter, fmt.Stringer) {
				s := &strings.Builder{}
				lw := NewLineWriter(s).LineBreak(" => ").Build()
				return lw, s
			},
			want: "wrap 2 => wrap => bail",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer, buf := tt.writer()
			tt.error.WriteError(writer)
			if got := buf.String(); got != tt.want {
				t.Errorf("Error.WriteError() = %v, want %v", got, tt.want)
			}
		})
	}
}
