package maleo

import (
	"context"
)

type detachedContext struct {
	context.Context
	inner context.Context
}

func (d detachedContext) Value(key any) any {
	return d.inner.Value(key)
}

// DetachedContext creates a context whose lifetime is detached from the input,
// but the values of the context is still reachable.
//
// DetachedContext is used by Maleo to send context.Context to Messengers that are not tied to the input lifetime.
func DetachedContext(ctx context.Context) context.Context {
	return detachedContext{Context: context.Background(), inner: ctx}
}
