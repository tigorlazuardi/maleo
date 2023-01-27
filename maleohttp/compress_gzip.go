package maleohttp

import (
	"bytes"
	"compress/gzip"
	"io"
)

var (
	_ Compressor       = (*GzipCompression)(nil)
	_ StreamCompressor = (*GzipCompression)(nil)
)

type GzipCompression struct {
	level int
}

// NewGzipCompression creates a new GzipCompression.
func NewGzipCompression() *GzipCompression {
	return NewGzipCompressionWithLevel(gzip.DefaultCompression)
}

// NewGzipCompressionWithLevel creates a new GzipCompression with specified compression level.
func NewGzipCompressionWithLevel(lvl int) *GzipCompression {
	return &GzipCompression{level: lvl}
}

// ContentEncoding implements maleohttp.ContentEncodingHint.
func (g GzipCompression) ContentEncoding() string {
	return "gzip"
}

// Compress implements maleohttp.Compressor.
func (g GzipCompression) Compress(b []byte) ([]byte, bool, error) {
	// 1500 is the max size of ethernet frame, 60 is the maximum range of TCP Header.
	//
	// The tradeoff between compression and cpu usage is not worth it if the size is less than MTU.
	//
	// Since the cost is the same: 1 IP packet.
	const minimumLength = 1500 - 60
	if len(b) < minimumLength {
		return b, false, nil
	}
	buf := new(bytes.Buffer)
	w, err := gzip.NewWriterLevel(buf, g.level)
	if err != nil {
		return b, false, err
	}
	_, err = w.Write(b)
	if err != nil {
		return b, false, err
	}
	_ = w.Close()
	c := make([]byte, buf.Len())
	// bytes.Buffer bytes method points to an array that will be reused by the pool.
	// So we need to copy the bytes to a new array.
	copy(c, buf.Bytes())
	return c, true, err
}

// StreamCompress implements maleohttp.StreamCompressor.
func (g GzipCompression) StreamCompress(contentType string, origin io.Reader) (io.Reader, bool) {
	// Gzip benefits heavily from text content. So we only compress text content.
	//
	// Images and other binary content often is already compressed, and compressing them again may actually increase the
	// size of the content.
	if !isHumanReadable(contentType) {
		return origin, false
	}
	pr, pw := io.Pipe()
	w, _ := gzip.NewWriterLevel(pw, gzip.BestCompression)
	go func() {
		_, err := io.Copy(w, origin)
		w.Close()
		_ = pw.CloseWithError(err)
	}()
	return pr, true
}
