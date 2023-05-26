package maleo

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrNil = errors.New("<nil>")

/*
ErrorBuilder is an interface to create customizable error.

ErrorBuilder by itself is not an error type. You have to call .Freeze() method to create proper Error type.
*/
type ErrorBuilder interface {
	// Code Sets the error code for this error.
	// This is used to identify the error type and how maleohttp will interact with this error.
	Code(i int) ErrorBuilder

	// Message Overrides the error message for this error.
	//
	// In built in implementation, If args are supplied, fmt.Sprintf will be called with s as base string.
	Message(s string, args ...any) ErrorBuilder

	// Error Sets the origin error for ErrorBuilder. Very unlikely to need to set this because Maleo.Wrap already uses the error.
	// But the api is available to set the origin error.
	Error(err error) ErrorBuilder

	// Context Sets additional data that will enrich how the error will look.
	//
	// The input is key-value format. The odd index will be used as key and the even index will be used as value.
	//
	// Key if not a string will be converted to string using fmt.Sprint.
	//
	// Example:
	//
	// 	maleo.Wrap(err).Code(400).Context("count", 123, "username", "kilua", "ranking", 5).Freeze()
	Context(ctx ...any) ErrorBuilder

	// Key Sets the key for this error. This is how the Messengers will use to identify if an error is the same as previous or not.
	//
	// Usually by not setting the key, The Messenger will generate their own.
	Key(key string, args ...any) ErrorBuilder

	// Caller Sets the caller for this error.
	Caller(c Caller) ErrorBuilder

	// Level Sets the level for this error.
	Level(lvl Level) ErrorBuilder

	// Time Sets the time for this error.
	Time(t time.Time) ErrorBuilder

	// Freeze this ErrorBuilder, preventing further mutations and set this ErrorBuilder into proper error.
	//
	// The returned Error is safe for multithreading usage because of its immutable nature.
	Freeze() Error

	// Log this error. Implicitly calls .Freeze() on this ErrorBuilder.
	Log(ctx context.Context) Error

	// Notify Notifies this error to Messengers. Implicitly calls .Freeze() on this ErrorBuilder.
	Notify(ctx context.Context, opts ...MessageOption) Error
}

type errorBuilder struct {
	code    int
	message string
	caller  Caller
	context []any
	key     string
	level   Level
	origin  error
	maleo   *Maleo
	time    time.Time
}

func (e *errorBuilder) Level(lvl Level) ErrorBuilder {
	e.level = lvl
	return e
}

func (e *errorBuilder) Caller(c Caller) ErrorBuilder {
	e.caller = c
	return e
}

func (e *errorBuilder) Code(i int) ErrorBuilder {
	e.code = i
	return e
}

func (e *errorBuilder) Error(err error) ErrorBuilder {
	if err == nil {
		err = ErrNil
	}
	e.origin = err
	return e
}

func (e *errorBuilder) Message(s string, args ...any) ErrorBuilder {
	if len(args) > 0 {
		e.message = fmt.Sprintf(s, args...)
	} else {
		e.message = s
	}
	return e
}

func (e *errorBuilder) Context(ctx ...any) ErrorBuilder {
	e.context = append(e.context, ctx...)
	return e
}

func (e *errorBuilder) Key(key string, args ...any) ErrorBuilder {
	if len(args) > 0 {
		e.key = fmt.Sprintf(key, args...)
	} else {
		e.key = key
	}
	return e
}

func (e *errorBuilder) Time(t time.Time) ErrorBuilder {
	e.time = t
	return e
}

func (e *errorBuilder) Freeze() Error {
	node := &ErrorNode{inner: e}
	if child, ok := e.origin.(*ErrorNode); ok {
		node.next = child
		child.prev = node
	}
	return node
}

func (e *errorBuilder) Log(ctx context.Context) Error {
	return e.Freeze().Log(ctx)
}

func (e *errorBuilder) Notify(ctx context.Context, opts ...MessageOption) Error {
	return e.Freeze().Notify(ctx, opts...)
}
