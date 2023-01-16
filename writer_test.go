package maleo

import (
	"bytes"
	"testing"
)

func TestNewLineWriter(t *testing.T) {
	out := new(bytes.Buffer)
	w := NewLineWriter(out).LineBreak("\n").Indent("  ").Prefix("[").Suffix("]").Build()
	if w == nil {
		t.Fatal("w is nil")
	}
	w.WritePrefix()
	n, err := w.WriteString("Hello")
	if err != nil {
		t.Fatal(err)
	}
	if n != 5 {
		t.Fatalf("n = %d, want 5", n)
	}
	w.WriteSuffix()
	if out.String() != "[Hello]" {
		t.Fatalf("out.String() = %q, want [Hello]", out.String())
	}
	if w.GetIndentation() != "  " {
		t.Fatalf("w.GetIndentation() = %q, want 2 spaces", w.GetIndentation())
	}
	if w.GetLineBreak() != "\n" {
		t.Fatalf("w.GetLineBreak() = %q, want \\n", w.GetLineBreak())
	}
	if w.GetPrefix() != "[" {
		t.Fatalf("w.GetPrefix() = %q, want [", w.GetPrefix())
	}
	if w.GetSuffix() != "]" {
		t.Fatalf("w.GetSuffix() = %q, want ]", w.GetSuffix())
	}
	out.Reset()
	w = NewLineWriter(w).LineBreak("\n").Indent("  ").Prefix("[").Suffix("]").Build()
	if w == nil {
		t.Fatal("w is nil")
	}
	w.WritePrefix()
	n, err = w.WriteString("Hello")
	if err != nil {
		t.Fatal(err)
	}
	if n != 5 {
		t.Fatalf("n = %d, want 5", n)
	}
	w.WriteSuffix()
	if out.String() != "[[Hello]]" {
		t.Fatalf("out.String() = %q, want [[Hello]]", out.String())
	}
	if w.GetIndentation() != "    " {
		t.Fatalf("w.GetIndentation() = %q, want 2 spaces", w.GetIndentation())
	}
	if w.GetLineBreak() != "\n\n" {
		t.Fatalf("w.GetLineBreak() = %q, want \\n", w.GetLineBreak())
	}
	if w.GetPrefix() != "[[" {
		t.Fatalf("w.GetPrefix() = %q, want [", w.GetPrefix())
	}
	if w.GetSuffix() != "]]" {
		t.Fatalf("w.GetSuffix() = %q, want ]", w.GetSuffix())
	}
}
