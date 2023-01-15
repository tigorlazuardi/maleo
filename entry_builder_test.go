package maleo

import (
	"context"
	"testing"
	"time"
)

func TestEntryBuilder_Caller(t *testing.T) {
	mal, _ := NewTestingMaleo()
	b := mal.NewEntry("foo")
	if b == nil {
		t.Fatal("Expected entry builder to be non-nil")
	}
	e := b.(*entryBuilder)
	if e.code != 0 {
		t.Errorf("Expected entry builder code to be 0, got %d", e.code)
	}
	if e.level != InfoLevel {
		t.Errorf("Expected entry builder level to be InfoLevel, got %s", e.level)
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
	m := newMockMessenger(1)
	mal.Register(m)
	b.Notify(context.Background())
	err := m.Wait(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %s", err)
	}
	if !m.called {
		t.Errorf("Expected messenger to be called")
	}
	l := newMockLogger()
	mal.SetLogger(l)
	b.Log(context.Background())
	if !l.called {
		t.Errorf("Expected logger to be called")
	}
	_ = b.Freeze().(EntryNode)
}
