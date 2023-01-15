package maleo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/kinbiko/jsonassert"
	"reflect"
	"testing"
	"time"
)

func TestEntryNode(t *testing.T) {
	mal, _ := NewTestingMaleo()
	now := time.Now()
	call := GetCaller(1)
	builder := mal.NewEntry("foo")
	builder.Code(600).
		Level(ErrorLevel).
		Key("foo").
		Time(now).
		Caller(call)

	node := builder.Freeze()
	if node == nil {
		t.Fatal("Expected entry node to be non-nil")
	}
	if node.Code() != 600 {
		t.Errorf("Expected entry node code to be 200, got %d", node.Code())
	}
	if node.HTTPCode() != 200 {
		t.Errorf("Expected entry node http code to be 200, got %d", node.HTTPCode())
	}
	if node.Service().String() != "test-test-test-test" {
		t.Errorf("Expected entry node service to be test-test-test-test, got %s", node.Service().String())
	}
	builder.Code(500)
	node = builder.Freeze()
	if node.HTTPCode() != 500 {
		t.Errorf("Expected entry node http code to be 500, got %d", node.HTTPCode())
	}
	builder.Code(1301)
	node = builder.Freeze()
	if node.HTTPCode() != 301 {
		t.Errorf("Expected entry node http code to be 301, got %d", node.HTTPCode())
	}
	if node.Level() != ErrorLevel {
		t.Errorf("Expected entry node level to be ErrorLevel, got %s", node.Level())
	}
	if node.Key() != "foo" {
		t.Errorf("Expected entry node key to be foo, got %s", node.Key())
	}
	if node.Message() != "foo" {
		t.Errorf("Expected entry node message to be foo, got %s", node.Message())
	}
	if !node.Time().Equal(now) {
		t.Errorf("Expected entry node time to be %s, got %s", now, node.Time())
	}
	if !reflect.DeepEqual(node.Caller(), call) {
		t.Errorf("Expected entry node caller to be %v, got %v", call, node.Caller())
	}

	builder.Context(1)
	node = builder.Freeze()
	if len(node.Context()) != 1 {
		t.Fatalf("Expected entry node context to be 1, got %d", len(node.Context()))
	}
	j := jsonassert.New(t)
	b, err := json.Marshal(node)
	if err != nil {
		t.Fatalf("Expected entry node to marshal to JSON without error, got %v", err)
	}
	j.Assertf(string(b), `
		{
			"time": "<<PRESENCE>>",
			"code": 1301,
			"message": "foo",
			"caller": "<<PRESENCE>>",
			"key": "foo",
			"level": "error",
			"service": {"name": "test", "environment": "test", "type": "test", "version": "v0.1.0-test"},
			"context": 1
		}`,
	)
	builder.Context(2)
	node = builder.Freeze()
	if len(node.Context()) != 2 {
		t.Fatalf("Expected entry node context to be 2, got %d", len(node.Context()))
	}
	b, err = json.Marshal(node)
	if err != nil {
		t.Fatalf("Expected entry node to marshal to JSON without error, got %v", err)
	}
	j.Assertf(string(b), `
		{
			"time": "<<PRESENCE>>",
			"code": 1301,
			"message": "foo",
			"caller": "<<PRESENCE>>",
			"key": "foo",
			"level": "error",
			"service": {"name": "test", "environment": "test", "type": "test", "version": "v0.1.0-test"},
			"context": [1, 2]
		}`)

	l := newMockLogger()
	mal.SetLogger(l)
	m := newMockMessenger(1)
	mal.Register(m)
	node.Notify(context.Background())
	err = mal.Wait(context.Background())
	if err != nil {
		t.Fatalf("Expected maleo to wait without error, got %v", err)
	}
	if !m.called {
		t.Error("Expected messenger to be called")
	}
	node.Log(context.Background())
	if !l.called {
		t.Error("Expected logger to be called")
	}

	if t.Failed() {
		out := new(bytes.Buffer)
		_ = json.Indent(out, b, "", "    ")
		fmt.Println(out.String())
	}
}
