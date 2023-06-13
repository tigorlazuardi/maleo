package maleo

import (
	"context"
	"fmt"
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

type contextKey struct {
	name string
}

func (co contextKey) String() string {
	return fmt.Sprintf("maleo.contextKey{name:%q}", co.name)
}

var contextKeyMaleo = contextKey{name: "maleo"}

// ContextWithMaleo creates a new context with Maleo instance attached.
//
// If there's an existing Maleo instance attached, it will be overwritten.
func ContextWithMaleo(parent context.Context, m *Maleo) context.Context {
	return context.WithValue(parent, contextKeyMaleo, m)
}

// MaleoFromContext retrieves Maleo instance from context.
//
// Returns nil if there is no Maleo instance attached.
func MaleoFromContext(ctx context.Context) *Maleo {
	m, _ := ctx.Value(contextKeyMaleo).(*Maleo)
	return m
}
