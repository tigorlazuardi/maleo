package maleo

import (
	"bytes"
	"context"
	"encoding/json"
	"time"
)

var _ Entry = (*EntryNode)(nil)

// EntryNode is the default implementation of Entry for Maleo.
type EntryNode struct {
	inner *entryBuilder
}

// MarshalJSON implements the json.Marshaler interface.
func (e EntryNode) MarshalJSON() ([]byte, error) {
	b := &bytes.Buffer{}
	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)
	err := enc.Encode(implJsonMarshaler{
		Time:    e.Time().Format(time.RFC3339),
		Code:    e.Code(),
		Message: e.Message(),
		Caller:  e.Caller(),
		Key:     e.Key(),
		Level:   e.Level().String(),
		Service: &e.inner.maleo.service,
		Context: e.Context(),
	})
	return b.Bytes(), err
}

// Code returns the original code of the type.
func (e EntryNode) Code() int {
	return e.inner.code
}

// Time returns the time of the entry.
func (e EntryNode) Time() time.Time {
	return e.inner.time
}

// Service returns the service name of the entry.
func (e EntryNode) Service() Service {
	return e.inner.maleo.service
}

// HTTPCode return HTTP Status Code for the type.
func (e EntryNode) HTTPCode() int {
	switch {
	case e.inner.code >= 200 && e.inner.code <= 599:
		return e.inner.code
	case e.inner.code > 999:
		code := e.inner.code % 1000
		if code >= 200 && code <= 599 {
			return code
		}
	}
	return 200
}

// Message returns the message.
func (e EntryNode) Message() string {
	return e.inner.message
}

// Key returns the key of the entry.
func (e EntryNode) Key() string {
	return e.inner.key
}

// Caller returns the maleo.Caller of the entry.
func (e EntryNode) Caller() Caller {
	return e.inner.caller
}

// Context returns the context of the entry.
func (e EntryNode) Context() []any {
	return e.inner.context
}

// Level Gets the level of this message.
func (e EntryNode) Level() Level {
	return e.inner.level
}

// Log logs the entry.
func (e EntryNode) Log(ctx context.Context) Entry {
	e.inner.maleo.Log(ctx, e)
	return e
}

// Notify sends the entry to the messengers.
func (e EntryNode) Notify(ctx context.Context, opts ...MessageOption) Entry {
	e.inner.maleo.Notify(ctx, e, opts...)
	return e
}
