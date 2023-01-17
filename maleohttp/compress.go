package maleohttp

import (
	"io"
)

type ContentEncodingHint interface {
	// ContentEncoding returns the value of the Content-Encoding header. If empty, the content-encoding header will be
	// set by the http framework.
	ContentEncoding() string
}

type Compressor interface {
	ContentEncodingHint
	// Compress the given bytes. If the compressed bytes are smaller than the original bytes, the compressed bytes will
	// be returned. Otherwise, the original bytes will be returned.
	//
	// If ok is true, the compressed bytes will be used and the Content-Encoding header will be set with the
	// value returned from ContentEncoding method. Otherwise, the original bytes will be used.
	//
	// If err is not nil, the original bytes will be used and the error will be logged by Maleo at warn level.
	Compress(b []byte) (compressed []byte, ok bool, err error)
}

type StreamCompressor interface {
	ContentEncodingHint
	// StreamCompress compresses the given reader and returns a new reader that will give the compressed data.
	//
	// If ok is true, the compressed bytes will be used and the Content-Encoding header will be set with the
	// value returned from ContentEncoding method. Otherwise, the original stream will be used.
	StreamCompress(contentType string, origin io.Reader) (io.Reader, bool)
}

var (
	_ Compressor       = (*NoCompression)(nil)
	_ StreamCompressor = (*NoCompression)(nil)
)

// NoCompression is a Compressor that does nothing. Basically it's an Uncompressed operation.
type NoCompression struct{}

func (n NoCompression) StreamCompress(_ string, origin io.Reader) (io.Reader, bool) {
	return origin, false
}
func (n NoCompression) ContentEncoding() string                 { return "" }
func (n NoCompression) Compress(b []byte) ([]byte, bool, error) { return b, false, nil }
