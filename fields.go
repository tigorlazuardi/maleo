package maleo

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Fields map[string]any

// F Alias to maleo.Fields.
type F = Fields

var (
	_ Summary       = (Fields)(nil)
	_ SummaryWriter = (Fields)(nil)
)

// Summary Returns a short summary of this type.
func (f Fields) Summary() string {
	s := &strings.Builder{}
	lw := NewLineWriter(s).LineBreak("\n").Build()
	f.WriteSummary(lw)
	return s.String()
}

// WriteSummary Writes the Summary() string to the writer instead of being allocated as value.
func (f Fields) WriteSummary(w LineWriter) {
	prefixLength := 0
	keys := make([]string, 0, len(f))
	for k := range f {
		if prefixLength < len(k) {
			prefixLength = len(k)
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		v := f[k]
		if i > 0 {
			w.WriteLineBreak()
		}
		i++

		w.WriteIndent()
		w.WritePrefix()
		_, _ = fmt.Fprintf(w, "%-*s: ", prefixLength, k)
		if v == nil {
			_, _ = w.WriteString("null")
			w.WriteSuffix()
			continue
		}
		switch v := v.(type) {
		case SummaryWriter:
			if _, ok := v.(Fields); ok {
				w.WriteLineBreak()
			}
			v.WriteSummary(NewLineWriter(w).Indent("  ").Build())
		case Summary:
			_, _ = w.WriteString(v.Summary())
		case fmt.Stringer:
			_, _ = w.WriteString(v.String())
		case json.RawMessage:
			if len(v) <= 32 {
				s := strconv.Quote(string(v))
				_, _ = w.WriteString(s)
			} else {
				_, _ = w.WriteString("[...]")
			}
		case []byte:
			if len(v) <= 32 {
				s := strconv.Quote(string(v))
				_, _ = w.WriteString(s)
			} else {
				_, _ = w.WriteString("[...]")
			}
		case string:
			if len(v) <= 32 {
				s := strconv.Quote(v)
				_, _ = w.WriteString(s)
			} else {
				_, _ = w.WriteString("[...]")
			}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, complex64, complex128:
			_, _ = fmt.Fprintf(w, "%v", v)
		default:
			_, _ = w.WriteString("[object]")
		}
		w.WriteSuffix()
	}
}
