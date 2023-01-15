package maleo

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

type TestingJSONLogger struct {
	buf *bytes.Buffer
	mu  sync.Mutex
}

// NewTestingMaleo returns a new Maleo instance with a TestingJSONLogger.
func NewTestingMaleo() (*Maleo, *TestingJSONLogger) {
	return NewTestingMaleoWithService(Service{
		Name:        "test",
		Environment: "test",
		Repository:  "",
		Branch:      "",
		Type:        "test",
		Version:     "v0.1.0-test",
	})
}

// NewTestingMaleoWithService returns a new Maleo instance with a TestingJSONLogger with a custom service metadata.
func NewTestingMaleoWithService(service Service) (*Maleo, *TestingJSONLogger) {
	logger := NewTestingJSONLogger()
	m := NewMaleo(service, Option.Init().Logger(logger))
	return m, logger
}

// NewTestingJSONLogger returns a very basic logger for testing purposes.
func NewTestingJSONLogger() *TestingJSONLogger {
	return &TestingJSONLogger{
		buf: new(bytes.Buffer),
	}
}

// Log implements maleo.Logger.
func (t *TestingJSONLogger) Log(ctx context.Context, entry Entry) {
	t.mu.Lock()
	defer t.mu.Unlock()

	err := json.NewEncoder(t.buf).Encode(entry)
	if err != nil {
		t.buf.WriteString(err.Error())
	}
}

// LogError implements maleo.Logger.
func (t *TestingJSONLogger) LogError(ctx context.Context, err Error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	errJson := json.NewEncoder(t.buf).Encode(err)
	if errJson != nil {
		t.buf.WriteString(errJson.Error())
	}
}

// Reset resets the buffer to be empty.
func (t *TestingJSONLogger) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.buf.Reset()
}

// String returns the accumulated bytes as string.
func (t *TestingJSONLogger) String() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.buf.String()
}

// Bytes returns the accumulated bytes.
func (t *TestingJSONLogger) Bytes() []byte {
	t.mu.Lock()
	defer t.mu.Unlock()
	cp := make([]byte, t.buf.Len())
	copy(cp, t.buf.Bytes())
	return cp
}

func (t *TestingJSONLogger) MarshalJSON() ([]byte, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]json.RawMessage, 0, 4)
	scanner := bufio.NewScanner(t.buf)
	for scanner.Scan() {
		out = append(out, json.RawMessage(scanner.Text()))
	}
	if len(out) == 1 {
		return out[0], nil
	}
	return json.Marshal(out)
}

func (t *TestingJSONLogger) PrettyPrint() {
	t.mu.Lock()
	defer t.mu.Unlock()
	var out bytes.Buffer
	scanner := bufio.NewScanner(t.buf)
	i := 0
	for scanner.Scan() {
		if i > 0 {
			out.WriteString("\n")
		}
		err := json.Indent(&out, scanner.Bytes(), "", "    ")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		i++
	}
	fmt.Println(out.String())
}
