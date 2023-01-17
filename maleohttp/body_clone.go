package maleohttp

import (
	"bytes"
	"io"
)

type ClonedBody interface {
	BufferedReader
	// CloneBytes returns a copy of the body bytes. Like String() but returns a copy of the bytes.
	CloneBytes() []byte
	// Truncated returns true if the body was truncated.
	Truncated() bool
	// Reader returns the reader that contains the body.
	//
	// Every call creates a new fresh cursor to the same underlying array. So multiple Readers created from this method
	// all have it's own read position.
	Reader() BufferedReader
}

type NoopCloneBody struct{}

func (NoopCloneBody) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (NoopCloneBody) Reader() BufferedReader {
	return NoopCloneBody{}
}

func (n NoopCloneBody) Bytes() []byte {
	return []byte{}
}

func (n NoopCloneBody) CloneBytes() []byte {
	return []byte{}
}

func (n NoopCloneBody) String() string {
	return ""
}

func (n NoopCloneBody) Len() int {
	return 0
}

func (n NoopCloneBody) Truncated() bool {
	return false
}

// BufferedReader is an extension around io.Reader that actually already have the values in memory and ready to be consumed.
type BufferedReader interface {
	io.Reader
	// String returns the body as a string. This is very often a copy operation of the bytes.
	String() string
	// Bytes returns the bytes of the body. This usually is not a copy operation, so the underlying array may be modified
	// somewhere else or in the future. Use String to ensure an immutable copy.
	Bytes() []byte
	// Len returns the number of bytes in the buffer.
	Len() int
}

type BufferedReadWriter interface {
	BufferedReader
	io.Writer
	// Reset resets the buffer to be empty, but it retains the underlying storage for use by future writes.
	Reset()
}

type noopReadWriter struct{ BufferedReader }

func (n2 noopReadWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (n2 noopReadWriter) Reset() {}

func wrapNoopWriter(reader BufferedReader) BufferedReadWriter {
	return noopReadWriter{reader}
}

// --------------------------------------------------------------------

var (
	_ ClonedBody         = (*bodyCloner)(nil)
	_ BufferedReadWriter = (*bodyCloner)(nil)
)

// bodyCloner is a wrapper around a Reader that clones the read body into a buffer.
type bodyCloner struct {
	io.ReadCloser
	clone BufferedReadWriter
	limit int
	cb    func(error)
}

func wrapBodyCloner(r io.Reader, limit int) *bodyCloner {
	if r == nil {
		return &bodyCloner{
			ReadCloser: io.NopCloser(NoopCloneBody{}),
			clone:      wrapNoopWriter(NoopCloneBody{}),
			limit:      0,
		}
	}
	var rc io.ReadCloser
	if rc2, ok := r.(io.ReadCloser); ok {
		rc = rc2
	} else {
		rc = io.NopCloser(r)
	}
	var cl BufferedReadWriter
	if buf, ok := r.(BufferedReader); ok {
		// Underlying type is already like bytes.Buffer, so no need for copy-like operations effectively. Since the data
		// is already in memory, we can just point to those arrays directly. We just have to make sure there's no write operations
		// to those underlying array to avoid double writing.
		cl = wrapNoopWriter(buf)
	} else {
		cl = &bytes.Buffer{}
	}
	return &bodyCloner{
		ReadCloser: rc,
		clone:      cl,
		limit:      limit,
	}
}

func (c *bodyCloner) onClose(cb func(error)) {
	c.cb = cb
}

func (c *bodyCloner) String() string {
	return c.clone.String()
}

func (c *bodyCloner) Bytes() []byte {
	return c.clone.Bytes()
}

func (c *bodyCloner) Len() int {
	return c.clone.Len()
}

func (c *bodyCloner) CloneBytes() []byte {
	res := make([]byte, c.clone.Len())
	copy(res, c.clone.Bytes())
	return res
}

func (c *bodyCloner) Truncated() bool {
	return c.limit > 0 && c.clone.Len() >= c.limit
}

func (c *bodyCloner) Reader() BufferedReader {
	// Create a new reader for a new fresh start whenever this method is called.
	return bytes.NewBuffer(c.clone.Bytes())
}

func (c *bodyCloner) Read(p []byte) (n int, err error) {
	n, err = c.ReadCloser.Read(p)
	// If we have a limit, we have to make sure we don't write more than the limit.
	if c.limit > 0 && c.clone.Len() >= c.limit {
		return n, err
	}
	// only write if limit is not 0.
	if n > 0 && c.limit != 0 {
		_, errWriteClone := c.clone.Write(p[:n])
		if err == nil {
			err = errWriteClone
		}
	}
	return n, err
}

func (c *bodyCloner) Close() error {
	err := c.ReadCloser.Close()
	if c.cb != nil {
		c.cb(err)
	}
	return err
}

func (c *bodyCloner) Write(p []byte) (n int, err error) {
	return c.clone.Write(p)
}

func (c *bodyCloner) Reset() {
	c.clone.Reset()
}
