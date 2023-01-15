package maleo

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestDetachedContext(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, "foo", "bar")
	dctx := DetachedContext(ctx)
	if dctx.Value("foo") != "bar" {
		t.Errorf("Expected detached context to have value bar, got %s", dctx.Value("foo"))
	}
	if reflect.TypeOf(dctx) != reflect.TypeOf(detachedContext{}) {
		t.Errorf("Expected detached context to be of type detachedContext, got %s", reflect.TypeOf(dctx))
	}
	if dctx.Done() != nil {
		t.Errorf("expected detached context done to be nil, got %#v", dctx.Done())
	}
}
