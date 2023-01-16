package bucket

import (
	"context"
	"io"
)

type File interface {
	Data() io.Reader
	Filename() string
	ContentType() string
	Read(p []byte) (n int, err error)
	Pretext() string
	Size() int
	Close() error
}

type implFile struct {
	data     io.Reader
	filename string
	mimetype string
	pretext  string
	size     int
}

func (f implFile) Data() io.Reader {
	return f.data
}

func (f implFile) Filename() string {
	return f.filename
}

func (f implFile) ContentType() string {
	return f.mimetype
}

func (f *implFile) Read(p []byte) (n int, err error) {
	return f.data.Read(p)
}

func (f implFile) Pretext() string {
	return f.pretext
}

func (f implFile) Size() int {
	return f.size
}

func (f *implFile) Close() error {
	if closer, ok := f.data.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// NewFile is a built-in constructor for File implementor.
func NewFile(data io.Reader, mimetype string, opts ...FileOption) File {
	var size int
	if lh, ok := data.(LengthHint); ok {
		size = lh.Len()
	}
	f := &implFile{
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

type Bucket interface {
	// Upload File(s) to the bucket.
	// If File.data implements io.Closer, the close method will be called after upload is done.
	// Whether the Upload operation is successful or not.
	//
	// The number of result will be the same as the number of files uploaded.
	Upload(ctx context.Context, files []File) []UploadResult
}

type FileOption interface {
	apply(*implFile)
}

type FileOptionFunc func(*implFile)

func (f FileOptionFunc) apply(file *implFile) {
	f(file)
}

func WithPretext(pretext string) FileOption {
	return FileOptionFunc(func(file *implFile) {
		file.pretext = pretext
	})
}

func WithFilesize(size int) FileOption {
	return FileOptionFunc(func(file *implFile) {
		file.size = size
	})
}

func WithFilename(filename string) FileOption {
	return FileOptionFunc(func(file *implFile) {
		file.filename = filename
	})
}
