package maleo

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestErrorBuilder(t *testing.T) {
	mal, _ := NewTestingMaleo()
	b := mal.Wrap(errors.New("foo"))
	if b == nil {
		t.Fatal("Expected entry builder to be non-nil")
	}
	e := b.(*errorBuilder)
	if e.code != 500 {
		t.Errorf("Expected entry builder code to be 500, got %d", e.code)
	}
	if e.level != ErrorLevel {
		t.Errorf("Expected entry builder level to be ErrorLevel, got %s", e.level)
	}
	if e.caller == nil {
		t.Errorf("Expected entry builder caller to be non-nil")
	}
	if e.message != "foo" {
		t.Errorf("Expected entry builder message to be foo, got %s", e.message)
	}
	b.Code(1000)
	if e.code != 1000 {
		t.Errorf("Expected entry builder code to be 1000, got %d", e.code)
	}
	b.Caller(GetCaller(1))
	if e.caller == nil {
		t.Errorf("Expected entry builder caller to be non-nil")
	}
	b.Context(F{"foo": "bar"})
	if len(e.context) != 1 {
		t.Errorf("Expected entry builder context to be 1, got %d", len(e.context))
	}
	b.Key("foo")
	if e.key != "foo" {
		t.Errorf("Expected entry builder key to be foo, got %s", e.key)
	}
	b.Message("baz")
	if e.message != "baz" {
		t.Errorf("Expected entry builder message to be baz, got %s", e.message)
	}
	b.Message("baz %s", "qux")
	if e.message != "baz qux" {
		t.Errorf("Expected entry builder message to be baz qux, got %s", e.message)
	}
	now := time.Now()
	b.Time(now)
	if !e.time.Equal(now) {
		t.Errorf("Expected entry builder time to be %s, got %s", now, e.time)
	}
	b.Level(ErrorLevel)
	if e.level != ErrorLevel {
		t.Errorf("Expected entry builder level to be ErrorLevel, got %s", e.level)
	}
	rre := errors.New("oof")
	b.Error(rre)
	if e.origin != rre {
		t.Errorf("Expected entry builder origin to be %s, got %s", rre, e.origin)
	}
	m := newMockMessenger(1)
	mal.Register(m)
	_ = b.Notify(context.Background())
	err := m.Wait(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}
	if !m.called {
		t.Errorf("Expected messenger to be called")
	}
	l := newMockLogger()
	mal.SetLogger(l)
	_ = b.Log(context.Background())
	if !l.called {
		t.Errorf("Expected logger to be called")
	}
	_ = b.Freeze().(Error)

	c := mal.Wrap(nil)
	d := c.(*errorBuilder)
	if d.origin == nil {
		t.Errorf("Expected nil error to be replaced with <nil>")
	}
	if d.message != "<nil>" {
		t.Errorf("Expected nil error to be replaced with <nil>")
	}
	gg := mal.BailFreeze("foo %s", "bar")
	h := mal.Wrap(gg)
	i := h.(*errorBuilder)
	if i.origin != gg {
		t.Errorf("Expected error to be wrapped")
	}
	if i.message != "foo bar" {
		t.Errorf("Expected error to be foo bar")
	}
	j := i.Freeze().(*ErrorNode)
	if j.next != gg {
		t.Errorf("Expected error.next to be ErrorNode")
	}
	k := mal.Wrap(j)
	n := k.(*errorBuilder)
	n.Error(nil)
	if n.origin != ErrNil {
		t.Errorf("Expected error to be <nil>")
	}
	n.Key("foo %s", "bar")
	if n.key != "foo bar" {
		t.Errorf("Expected key to be foo bar")
	}
}
