package maleo

import (
	"context"
	"encoding/json"
	"testing"
)

func TestNewTestingJSONLogger(t *testing.T) {
	mal, log := NewTestingMaleo()
	mal.NewEntry("").Log(context.Background())
	if len(log.Bytes()) == 0 {
		t.Fatal("Expected log to have bytes")
	}
	if len(log.String()) == 0 {
		t.Fatal("Expected log to have string")
	}
	log.PrettyPrint()
	_ = mal.BailFreeze("foo").Log(context.Background())
	marshal, err := json.Marshal(log)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(marshal) == 0 {
		t.Fatal("Expected marshal to have bytes")
	}
	log.Reset()
	if len(log.Bytes()) != 0 {
		t.Fatal("Expected log to be empty")
	}
}
