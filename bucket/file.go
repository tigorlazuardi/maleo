package bucket

import (
	"io"
)

// --8<-- [start:file]

type File interface {
	Data() io.Reader
	Filename() string
	ContentType() string
	Read(p []byte) (n int, err error)
	Pretext() string
	Size() int
	Close() error
}

// --8<-- [end:file]

type file struct {
	data     io.Reader
	filename string
	mimetype string
	pretext  string
	size     int
}

func (f *file) Data() io.Reader {
	return f.data
}

func (f *file) Filename() string {
	return f.filename
}

func (f *file) ContentType() string {
	return f.mimetype
}

func (f *file) Read(p []byte) (n int, err error) {
	return f.data.Read(p)
}

func (f *file) Pretext() string {
	return f.pretext
}

func (f *file) Size() int {
	return f.size
}

func (f *file) Close() error {
	if closer, ok := f.data.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// NewFile is a built-in constructor for File implementor.
//
// if data implements io.Closer, it will be closed when the file is uploaded by the Bucket.
//
// If WithFilename option is not provided, a randomly generated snowflake ID will be used.
//
// if data is a reader that have method .Len() int it will be used to set the file size. (e.g. bytes.Buffer).
// Otherwise, the size will be set to -1. You may use WithFilesize to set the size manually.
func NewFile(data io.Reader, mimetype string, opts ...FileOption) File {
	size := -1
	if lh, ok := data.(LengthHint); ok {
		size = lh.Len()
	}
	f := &file{
		data:     data,
		filename: snowflakeNode.Generate().String(),
		mimetype: mimetype,
		size:     size,
	}
	for _, opt := range opts {
		opt.apply(f)
	}
	return f
}

type UploadResult struct {
	// The URL of the uploaded file, if successful.
	URL string
	// The file instance used to upload the file.
	// The body of this file may have already been garbage collected.
	// So do not consume this file content again and only use the remaining metadata.
	File File
	// If Error is not nil, the upload is considered failed.
	Error error
}

type FileOption interface {
	apply(*file)
}

type FileOptionFunc func(*file)

func (f FileOptionFunc) apply(file *file) {
	f(file)
}

// WithPretext sets a description for the file.
func WithPretext(pretext string) FileOption {
	return FileOptionFunc(func(file *file) {
		file.pretext = pretext
	})
}

// WithFilesize sets the size of the file.
func WithFilesize(size int) FileOption {
	return FileOptionFunc(func(file *file) {
		file.size = size
	})
}

// WithFilename sets the filename of the file. Default is a randomly generated snowflake ID.
func WithFilename(filename string) FileOption {
	return FileOptionFunc(func(file *file) {
		file.filename = filename
	})
}
