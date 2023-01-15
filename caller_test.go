package maleo

import (
	"bytes"
	"strings"
	"testing"
)

func TestCaller(t *testing.T) {
	c := GetCaller(1)
	if c == nil {
		t.Fatal("Expected caller to be non-nil")
	}
	if c.Name() != "github.com/tigorlazuardi/maleo.TestCaller" {
		t.Errorf("Expected caller name to be github.com/tigorlazuardi/maleo.TestCaller, got %s", c.Name())
	}
	if c.ShortName() != "maleo.TestCaller" {
		t.Errorf("Expected caller short name to be maleo.TestCaller, got %s", c.ShortName())
	}
	if !strings.Contains(c.File(), "caller_test.go") {
		t.Errorf("Expected caller file to be caller_test.go, got %s", c.File())
	}
	if c.PC() == 0 {
		t.Errorf("Expected caller pc to be non-zero")
	}
	if !strings.Contains(c.FormatAsKey(), "caller_test.go_") {
		t.Errorf("Expected caller format as key to be caller_test.go_, got %s", c.FormatAsKey())
	}
	if c.Line() <= 2 {
		t.Errorf("expected caller line to be greater than 2, got %d", c.Line())
	}
	if c.Function() == nil {
		t.Errorf("Expected caller function to be non-nil")
	}
	if !strings.Contains(c.String(), "caller_test.go:") {
		t.Errorf("Expected caller string to be caller_test.go:, got %s", c.String())
	}
	if !strings.Contains(c.ShortSource(), "caller_test.go") {
		t.Errorf("Expected caller short source to be caller_test.go, got %s", c.ShortSource())
	}
	if c.Depth() != 1 {
		t.Errorf("Expected caller depth to be 1, got %d", c.Depth())
	}
	cal := c.(*caller)
	b, err := cal.MarshalJSON()
	if err != nil {
		t.Errorf("Expected caller marshal json to be nil, got %s", err)
	}
	if !bytes.Contains(b, []byte("caller_test.go:")) {
		t.Errorf("Expected caller marshal json to be caller_test.go:, got %s", b)
	}
}
