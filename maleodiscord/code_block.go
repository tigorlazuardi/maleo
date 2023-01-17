package maleodiscord

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"

	"github.com/tigorlazuardi/maleo"
)

type CodeBlockBuilder interface {
	Build(w io.Writer, value []any) error
	BuildError(w io.Writer, err error) error
}

type JSONCodeBlockBuilder struct{}

type valueMarshaler []any

var _ maleo.CodeBlockJSONMarshaler = (valueMarshaler)(nil)

func (v valueMarshaler) CodeBlockJSON() ([]byte, error) {
	const indent = "   "
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if len(v) == 1 {
		enc.SetIndent("", indent)
		err := enc.Encode(v[0])
		return buf.Bytes(), err
	}
	j := make([]json.RawMessage, 0, len(v))
	for _, value := range v {
		buf.Reset()
		if vm, ok := value.(maleo.CodeBlockJSONMarshaler); ok {
			raw, err := vm.CodeBlockJSON()
			if err != nil {
				return nil, err
			}
			j = append(j, raw)
			continue
		}
		err := enc.Encode(value)
		if err != nil {
			return nil, err
		}
		j = append(j, json.RawMessage(buf.String()))
	}
	buf.Reset()
	enc.SetIndent("", indent)
	_ = enc.Encode(j)
	return buf.Bytes(), nil
}

func (J JSONCodeBlockBuilder) Build(w io.Writer, value []any) error {
	_, err := io.WriteString(w, "```json\n")
	if err != nil {
		return err
	}
	defer func(w io.Writer, s string) {
		_, _ = io.WriteString(w, s)
	}(w, "```")
	b, err := valueMarshaler(value).CodeBlockJSON()
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (J JSONCodeBlockBuilder) BuildError(w io.Writer, e error) error {
	_, err := io.WriteString(w, "```json\n")
	if err != nil {
		return err
	}
	defer func(w io.Writer, s string) {
		_, _ = io.WriteString(w, s)
	}(w, "```")
	if e, ok := e.(maleo.CodeBlockJSONMarshaler); ok {
		b, err := e.CodeBlockJSON()
		if err != nil {
			return err
		}
		_, err = w.Write(b)
		if err != nil {
			return err
		}
		return nil
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "   ")
	enc.SetEscapeHTML(false)
	return enc.Encode(richJSONError{e})
}

type richJSONError struct {
	err error
}

func (r richJSONError) MarshalJSON() ([]byte, error) {
	if r.err == nil {
		return []byte(`{"error":null}`), nil
	}
	if jm, ok := r.err.(json.Marshaler); ok {
		return jm.MarshalJSON()
	}
	b, err := json.Marshal(r.err)
	if err != nil {
		return nil, err
	}
	w := new(bytes.Buffer)
	w.WriteString(`{"error":`)
	switch {
	case len(b) < 2,
		b[0] == '{' && b[1] == '}',
		b[0] == '[' && b[1] == ']':

		w.WriteString(strconv.Quote(r.err.Error()))
		w.WriteRune('}')
		return w.Bytes(), nil
	}
	w.WriteString(`{"summary":"`)
	w.WriteString(r.err.Error())
	w.WriteString(`","details":`)
	w.Write(b)
	w.WriteString("}}")
	return w.Bytes(), nil
}
