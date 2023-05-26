package maleo

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

const codeBlockIndent = "   "

// ErrorNode is the implementation of the Error interface.
type ErrorNode struct {
	inner *errorBuilder
	prev  *ErrorNode
	next  *ErrorNode
}

// sorted keys are rather important for human reads. Especially the Context and Error should always be at the last marshaled keys.
// as they contain the most amount of data and information, and thus shadows other values at a glance.
//
// arguably this is simpler to be done than implementing json.Marshaler interface and doing it manually, key by key
// without resorting to other libraries.
type implJsonMarshaler struct {
	Time    string   `json:"time,omitempty"`
	Code    int      `json:"code,omitempty"`
	Message string   `json:"message,omitempty"`
	Caller  Caller   `json:"caller,omitempty"`
	Key     string   `json:"key,omitempty"`
	Level   string   `json:"level,omitempty"`
	Service *Service `json:"service,omitempty"`
	Context any      `json:"context,omitempty"`
	Error   error    `json:"error,omitempty"`
}

func newImplJSONMarshaler(e Error, next error, ctx any, service *Service) implJsonMarshaler {
	return implJsonMarshaler{
		Time:    e.Time().Format(time.RFC3339),
		Code:    e.Code(),
		Message: e.Message(),
		Caller:  e.Caller(),
		Key:     e.Key(),
		Level:   e.Level().String(),
		Context: ctx,
		Error:   richJsonError{next},
		Service: service,
	}
}

type marshalFlag uint8

func (m marshalFlag) Has(f marshalFlag) bool {
	return m&f == f
}

func (m *marshalFlag) Set(f marshalFlag) {
	*m |= f
}

func (m *marshalFlag) Unset(f marshalFlag) {
	*m &= ^f
}

const (
	marshalSkipCode marshalFlag = 1 << iota
	marshalSkipMessage
	marshalSkipLevel
	marshalSkipCaller
	marshalSkipContext
	marshalSkipTime
	marshalSkipService
	marshalSkipAll = marshalSkipCode +
		marshalSkipMessage +
		marshalSkipLevel +
		marshalSkipTime +
		marshalSkipContext +
		marshalSkipCaller +
		marshalSkipService
)

func (e *ErrorNode) createMarshalJSONFlag() marshalFlag {
	var m marshalFlag
	// The logic below for condition flow:
	//
	// if the next error is not an ErrorNode, denoted by e.next == nil, we will test against the error implements
	// Error interface, and deduplicate the fields in this current node. But only when the previous node is also an
	// ErrorNode.
	//
	// Unlike in CodeBlock for human read first, where the innermost error is the most important.
	//
	// the Logic for MarshalJSON is aimed towards machine and log parsers.
	//
	// The outermost error is the most important fields for indexing, and thus we will not skip any fields.
	//
	// However, any nested error with duplicate values will be just a waste of space and bandwidth, so we will skip them.
	if e.prev == nil {
		return m
	}
	other, ok := e.inner.origin.(Error)
	if e.prev == nil && e.next == nil && !ok {
		return m
	}
	prev, current := e.prev, e
	if prev.Code() == current.Code() {
		m.Set(marshalSkipCode)
	}
	if prev.Level() == current.Level() {
		m.Set(marshalSkipLevel)
	}
	// e.next != nil because we don't want message skipped when it's the last error in the chain.
	if prev.Message() == current.Message() && e.next != nil {
		m.Set(marshalSkipMessage)
	}
	if len(current.Context()) == 0 {
		m.Set(marshalSkipContext)
	}
	if prev.Time().Sub(current.Time()) < time.Second {
		m.Set(marshalSkipTime)
	}
	if prev.inner.maleo.service == current.inner.maleo.service {
		m.Set(marshalSkipService)
	}
	if ok {
		m |= e.deduplicateAgainstOtherError(other)
	}
	if m.Has(marshalSkipCode) &&
		m.Has(marshalSkipMessage) &&
		m.Has(marshalSkipLevel) &&
		m.Has(marshalSkipContext) &&
		m.Has(marshalSkipTime) &&
		m.Has(marshalSkipService) {
		m.Set(marshalSkipCaller)
	}
	return m
}

func (e *ErrorNode) createPayload(m marshalFlag) *implJsonMarshaler {
	ctx := func() any {
		if len(e.inner.context) == 0 {
			return nil
		}
		if len(e.inner.context) == 1 {
			return e.inner.context[0]
		}
		return e.inner.context
	}()
	var next error
	if e.next != nil {
		next = e.next
	} else {
		next = e.inner.origin
	}
	marshalAble := newImplJSONMarshaler(e, next, ctx, &e.inner.maleo.service)

	if m.Has(marshalSkipCode) {
		marshalAble.Code = 0
	}
	if m.Has(marshalSkipMessage) {
		marshalAble.Message = ""
	}
	if m.Has(marshalSkipLevel) {
		marshalAble.Level = ""
	}
	if m.Has(marshalSkipTime) {
		marshalAble.Time = ""
	}
	if m.Has(marshalSkipCaller) {
		marshalAble.Caller = nil
	}
	if m.Has(marshalSkipService) {
		marshalAble.Service = nil
	}
	return &marshalAble
}

func (e *ErrorNode) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte("null"), nil
	}
	m := e.createMarshalJSONFlag()
	b := &bytes.Buffer{}
	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)
	if m.Has(marshalSkipAll) {
		err := enc.Encode(richJsonError{e.inner.origin})
		return b.Bytes(), err
	}
	err := enc.Encode(e.createPayload(m))
	return b.Bytes(), err
}

func (e *ErrorNode) Error() string {
	s := &strings.Builder{}
	lw := NewLineWriter(s).LineBreak(": ").Build()
	e.WriteError(lw)
	return s.String()
}

// WriteError Writes the error.Error to the writer instead of being allocated as value.
func (e *ErrorNode) WriteError(w LineWriter) {
	w.WriteIndent()
	msg := e.inner.message
	if e.inner.origin == nil {
		// Account for empty string message after wrapping nil error.
		if len(msg) > 0 {
			w.WritePrefix()
			_, _ = w.WriteString(msg)
			w.WriteSuffix()
			w.WriteLineBreak()
		}
		w.WritePrefix()
		_, _ = w.WriteString("[nil]")
		w.WriteSuffix()
		return
	}

	writeInner := func(linebreak bool) {
		if ew, ok := e.inner.origin.(ErrorWriter); ok {
			if linebreak {
				w.WriteLineBreak()
			}
			ew.WriteError(w)
		} else {
			errMsg := e.inner.origin.Error()
			if errMsg != msg {
				w.WriteLineBreak()
				w.WritePrefix()
				_, _ = w.WriteString(errMsg)
				w.WriteSuffix()
			}
		}
	}

	var innerMessage string
	if mh, ok := e.inner.origin.(MessageHint); ok {
		innerMessage = mh.Message()
	}

	// Skip writing duplicate or empty messages.
	if msg == innerMessage || len(msg) == 0 {
		writeInner(false)
		return
	}

	w.WritePrefix()
	_, _ = w.WriteString(msg)
	w.WriteSuffix()
	writeInner(true)
}

// Code Gets the original code of the type.
func (e *ErrorNode) Code() int {
	return e.inner.code
}

// HTTPCode Gets HTTP Status Code for the type.
func (e *ErrorNode) HTTPCode() int {
	switch {
	case e.inner.code >= 200 && e.inner.code <= 599:
		return e.inner.code
	case e.inner.code > 999:
		code := e.inner.code % 1000
		if code >= 200 && code <= 599 {
			return code
		}
	}
	return 500
}

// Message Gets the Message of the type.
func (e *ErrorNode) Message() string {
	return e.inner.message
}

// Caller Gets the caller of this type.
func (e *ErrorNode) Caller() Caller {
	return e.inner.caller
}

// Context Gets the context of this type.
func (e *ErrorNode) Context() []any {
	return e.inner.context
}

func (e *ErrorNode) Level() Level {
	return e.inner.level
}

func (e *ErrorNode) Time() time.Time {
	return e.inner.time
}

func (e *ErrorNode) Key() string {
	return e.inner.key
}

func (e *ErrorNode) Service() Service {
	return e.inner.maleo.service
}

// Unwrap Returns the error that is wrapped by this error. To be used by errors.Is and errors.As functions from errors library.
func (e *ErrorNode) Unwrap() error {
	if e.next != nil {
		return e.next
	}
	return e.inner.origin
}

// Log this error.
func (e *ErrorNode) Log(ctx context.Context) Error {
	e.inner.maleo.LogError(ctx, e)
	return e
}

// Notify this error to Messengers.
func (e *ErrorNode) Notify(ctx context.Context, opts ...MessageOption) Error {
	e.inner.maleo.NotifyError(ctx, e, opts...)
	return e
}

// richJsonError is a special kind of error that tries to prevent information loss when marshaling to json.
type richJsonError struct {
	error
}

func (r richJsonError) MarshalJSON() ([]byte, error) {
	if r.error == nil {
		return []byte("null"), nil
	}
	// if the error supports json.Marshaler we use it directly.
	// this is because we can assume that the error have special marshaling needs for specific output.
	//
	// E.G. to prevent unnecessary "summary" keys when the origin error is already is a maleo.Error type.
	if e, ok := r.error.(json.Marshaler); ok { //nolint
		return e.MarshalJSON()
	}
	b := &bytes.Buffer{}
	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)

	err := enc.Encode(r.error)
	if err != nil {
		_ = enc.Encode(r.error.Error())
		return b.Bytes(), nil
	}

	summary := r.error.Error()

	// We now handle empty errors after encoding.
	//
	// 3 because it also includes newline after brackets or quotes.
	//
	// The logic below looks up for: ""\n, {}\n, []\n
	if b.Len() == 3 && b.Bytes()[2] == '\n' {
		v := b.Bytes()
		switch {
		case v[0] == '"', v[0] == '{', v[0] == '[':
			b.Reset()
			err = enc.Encode(map[string]string{"summary": summary})
			return b.Bytes(), err
		}
	}

	content := b.String()
	b.Reset()
	err = enc.Encode(map[string]json.RawMessage{
		"summary": json.RawMessage(strconv.Quote(summary)),
		"details": json.RawMessage(content),
	})
	return b.Bytes(), err
}
