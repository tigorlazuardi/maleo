package maleo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type CodeBlockJSONMarshaler interface {
	CodeBlockJSON() ([]byte, error)
}

type cbJson struct {
	inner error
}

func (c cbJson) Error() string {
	return c.inner.Error()
}

func (c cbJson) CodeBlockJSON() ([]byte, error) {
	return c.MarshalJSON()
}

func (c cbJson) MarshalJSON() ([]byte, error) {
	if c.inner == nil {
		return []byte("null"), nil
	}
	if cb, ok := c.inner.(CodeBlockJSONMarshaler); ok && cb != nil {
		return cb.CodeBlockJSON()
	}
	b := &bytes.Buffer{}
	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", codeBlockIndent)
	err := enc.Encode(richJsonError{c.inner})
	return b.Bytes(), err
}

func (e *ErrorNode) CodeBlockJSON() ([]byte, error) {
	if e == nil {
		return []byte("null"), nil
	}
	m := e.createCodeBlockMarshalFlag()
	b := &bytes.Buffer{}
	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", codeBlockIndent)
	// Check if current ErrorNode needs to be skipped.
	if m.Has(marshalSkipAll) {
		origin := e.inner.origin
		if cbJson, ok := origin.(CodeBlockJSONMarshaler); ok {
			return cbJson.CodeBlockJSON()
		}
		err := enc.Encode(richJsonError{origin})
		return b.Bytes(), err
	}
	err := enc.Encode(e.createCodeBlockPayload(m))
	return bytes.TrimSpace(b.Bytes()), err
}

// createCodeBlockMarshalFlag creates a flag that skips the fields that have the same value as the parent *ErrorNode.
func (e *ErrorNode) createCodeBlockMarshalFlag() marshalFlag {
	var m marshalFlag
	origin, ok := e.inner.origin.(Error)
	if !ok {
		return m
	}
	if origin.Code() == e.Code() {
		m.Set(marshalSkipCode)
	}
	if origin.Message() == e.Message() {
		m.Set(marshalSkipMessage)
	}
	if origin.Level() == e.Level() {
		m.Set(marshalSkipLevel)
	}
	if len(origin.Context()) == 0 {
		m.Set(marshalSkipContext)
	}
	if origin.Time().Sub(e.Time()) < time.Second {
		m.Set(marshalSkipTime)
	}
	if m.Has(marshalSkipCode) &&
		m.Has(marshalSkipMessage) &&
		m.Has(marshalSkipLevel) &&
		m.Has(marshalSkipContext) {
		m.Set(marshalSkipCaller)
	}
	if origin.Service() == e.Service() {
		m.Set(marshalSkipService)
	}
	return m
}

func (e *ErrorNode) deduplicateAgainstOtherError(other Error) marshalFlag {
	var m marshalFlag
	if e.Code() == other.Code() {
		m.Set(marshalSkipCode)
	}
	if e.Message() == other.Message() {
		m.Set(marshalSkipMessage)
	}
	if e.Level() == other.Level() {
		m.Set(marshalSkipLevel)
	}
	if len(other.Context()) == 0 {
		m.Set(marshalSkipContext)
	}
	if e.Time().Sub(other.Time()) < time.Second {
		m.Set(marshalSkipTime)
	}
	if e.inner.maleo.service == other.Service() {
		m.Set(marshalSkipService)
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

func toMap(v []any) map[string]any {
	var (
		key   string
		value any
		out   = make(map[string]any, len(v)/2+1)
	)
	for i := 0; i < len(v); i++ {
		if i%2 == 0 {
			keyAssert, ok := v[i].(string)
			if ok {
				key = keyAssert
			} else {
				key = fmt.Sprint(v[i])
			}
		} else {
			value = v[i]
		}
		if key != "" && value != nil {
			out[key] = value
			key = ""
			value = nil
		}
	}
	return out
}

func (e *ErrorNode) createCodeBlockPayload(m marshalFlag) *implJsonMarshaler {
	ctx := func() any {
		if len(e.inner.context) == 0 {
			return nil
		}
		if len(e.inner.context) == 1 {
			return e.inner.context[0]
		}
		return toMap(e.inner.context)
	}()
	var next error
	if e.next != nil {
		next = e.next
	} else {
		next = e.inner.origin
	}
	marshalAble := implJsonMarshaler{
		Time:    e.Time().Format(time.RFC3339),
		Code:    e.Code(),
		Message: e.Message(),
		Caller:  e.Caller(),
		Key:     e.Key(),
		Level:   e.Level().String(),
		Context: ctx,
		Error:   cbJson{next},
		Service: &e.inner.maleo.service,
	}

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
