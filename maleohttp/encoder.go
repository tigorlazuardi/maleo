package maleohttp

import (
	"bytes"
	"encoding/json"
	"sync"
)

type ContentTypeHint interface {
	// ContentType Returns the content type of the encoded data.
	ContentType() string
}

type Encoder interface {
	ContentTypeHint
	// Encode the input into a byte array.
	Encode(input any) ([]byte, error)
}

type JSONEncoder struct {
	htmlEscape bool
	indent     string
	prefix     string
	pool       *sync.Pool
}

// SetHtmlEscape Sets the htmlEscape flag. If set to true, HTML characters will be escaped.
// Useful if you plan to embed HTML in a JSON field.
func (j *JSONEncoder) SetHtmlEscape(htmlEscape bool) {
	j.htmlEscape = htmlEscape
}

// SetIndent Sets the indent for every level line. Used for pretty print JSON. Empty value disable indentation.
func (j *JSONEncoder) SetIndent(indent string) {
	j.indent = indent
}

// SetPrefix Sets the prefix of the output for every line. Empty value disable prefix.
func (j *JSONEncoder) SetPrefix(prefix string) {
	j.prefix = prefix
}

// NewJSONEncoder Creates a new JSONEncoder.
func NewJSONEncoder() *JSONEncoder {
	return &JSONEncoder{
		pool: &sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}
}

func (j *JSONEncoder) ContentType() string {
	return "application/json"
}

func (j *JSONEncoder) Encode(input any) ([]byte, error) {
	buf := j.pool.Get().(*bytes.Buffer) // nolint
	buf.Reset()
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(j.htmlEscape)
	enc.SetIndent(j.prefix, j.indent)
	err := enc.Encode(input)
	b := buf.Bytes()
	c := make([]byte, len(b))
	copy(c, b)
	j.pool.Put(buf)
	return c, err
}
