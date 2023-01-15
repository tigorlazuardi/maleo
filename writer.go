package maleo

import "io"

type LineWriter interface {
	io.Writer
	io.StringWriter
	// WriteLineBreak Writes a predetermined new line character(s) to the writer.
	WriteLineBreak()
	// WritePrefix Writes a predetermined prefix to the writer.
	WritePrefix()
	// WriteSuffix Writes a predetermined suffix to the writer.
	WriteSuffix()
	// WriteIndent Writes Indentation characters.
	WriteIndent()
	// GetLineBreak Returns the pre-determined line break characters.
	GetLineBreak() string
	// GetPrefix Returns the pre-determined prefix.
	GetPrefix() string
	// GetSuffix Returns the pre-determined suffix.
	GetSuffix() string
	// GetIndentation Returns the pre-determined indentation.
	GetIndentation() string
}

type LineWriterBuilder struct {
	writer  io.Writer
	indent  string
	newLine string
	prefix  string
	suffix  string
}

// Indent Sets the Indentation.
func (builder *LineWriterBuilder) Indent(s string) *LineWriterBuilder {
	builder.indent = s
	return builder
}

// LineBreak Sets the Linebreak character(s).
func (builder *LineWriterBuilder) LineBreak(s string) *LineWriterBuilder {
	builder.newLine = s
	return builder
}

// Prefix Sets the Prefix.
func (builder *LineWriterBuilder) Prefix(s string) *LineWriterBuilder {
	builder.prefix = s
	return builder
}

// Suffix Sets the Suffix.
func (builder *LineWriterBuilder) Suffix(s string) *LineWriterBuilder {
	builder.suffix = s
	return builder
}

// Build Turn this writer into proper LineWriter.
func (builder *LineWriterBuilder) Build() LineWriter {
	return &lineWriter{
		Writer:  builder.writer,
		newLine: builder.newLine,
		prefix:  builder.prefix,
		suffix:  builder.suffix,
		indent:  builder.indent,
	}
}

// NewLineWriter Creates a new LineWriterBuilder. You have to call .Build() to actually use LineWriter.
func NewLineWriter(writer io.Writer) *LineWriterBuilder {
	return &LineWriterBuilder{
		writer: writer,
	}
}

var _ LineWriter = (*lineWriter)(nil)

type lineWriter struct {
	io.Writer
	newLine string
	prefix  string
	suffix  string
	indent  string
}

func (l *lineWriter) WriteString(s string) (n int, err error) {
	return io.WriteString(l.Writer, s)
}

// WriteIndent Writes Indentation characters.
func (l *lineWriter) WriteIndent() {
	if lw, ok := l.Writer.(LineWriter); ok {
		lw.WriteIndent()
	}
	_, _ = io.WriteString(l.Writer, l.indent)
}

func (l *lineWriter) GetIndentation() string {
	if lw, ok := l.Writer.(LineWriter); ok {
		return l.indent + lw.GetIndentation()
	}
	return l.indent
}

func (l *lineWriter) WriteLineBreak() {
	if lw, ok := l.Writer.(LineWriter); ok {
		lw.WriteLineBreak()
	}
	_, _ = io.WriteString(l.Writer, l.newLine)
}

func (l *lineWriter) WritePrefix() {
	if lw, ok := l.Writer.(LineWriter); ok {
		lw.WritePrefix()
	}
	_, _ = io.WriteString(l.Writer, l.prefix)
}

func (l *lineWriter) WriteSuffix() {
	if lw, ok := l.Writer.(LineWriter); ok {
		lw.WriteSuffix()
	}
	_, _ = io.WriteString(l.Writer, l.suffix)
}

func (l *lineWriter) GetLineBreak() string {
	if lw, ok := l.Writer.(LineWriter); ok {
		return lw.GetLineBreak() + l.newLine
	}
	return l.newLine
}

func (l *lineWriter) GetPrefix() string {
	if lw, ok := l.Writer.(LineWriter); ok {
		return lw.GetPrefix() + l.prefix
	}
	return l.prefix
}

func (l *lineWriter) GetSuffix() string {
	if lw, ok := l.Writer.(LineWriter); ok {
		return lw.GetSuffix() + l.suffix
	}
	return l.suffix
}

type Display interface {
	// Display returns a human-readable and rich with information for the implementer.
	Display() string
}

type DisplayWriter interface {
	// WriteDisplay Writes the Display() string to the writer instead of being allocated as value.
	WriteDisplay(w LineWriter)
}

type ErrorWriter interface {
	// WriteError Writes the error.Error to the writer instead of being allocated as value.
	WriteError(w LineWriter)
}

type Summary interface {
	// Summary Returns a short summary of the implementer.
	Summary() string
}

type SummaryWriter interface {
	// WriteSummary Writes the Summary() string to the writer instead of being allocated as value.
	WriteSummary(w LineWriter)
}
