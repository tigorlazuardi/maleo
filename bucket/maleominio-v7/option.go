package maleominio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/tigorlazuardi/maleo/bucket"
)

type Option struct {
	makeBucketOption minio.MakeBucketOptions
	putObjectOptions PutObjectOptionFunc
	prefix           fmt.Stringer
}

type WrapOption interface {
	apply(*Option)
}

type wrapOptionFunc func(*Option)

func (f wrapOptionFunc) apply(o *Option) {
	f(o)
}

// WithMakeBucketOption sets make bucket option for Minio if bucket does not exist.
func WithMakeBucketOption(opt minio.MakeBucketOptions) WrapOption {
	return wrapOptionFunc(func(o *Option) {
		o.makeBucketOption = opt
	})
}

type PutObjectOptionFunc = func(ctx context.Context, file bucket.File) minio.PutObjectOptions

// WithPutObjectOption sets put object option for Minio.
func WithPutObjectOption(f PutObjectOptionFunc) WrapOption {
	return wrapOptionFunc(func(o *Option) {
		o.putObjectOptions = f
	})
}

type StringerFunc func() string

func (s StringerFunc) String() string {
	return s()
}

// WithFilePrefix sets static prefix for file name.
func WithFilePrefix(s string) WrapOption {
	return WithFilePrefixStringer(StringerFunc(func() string {
		return s
	}))
}

// WithFilePrefixStringer sets dynamic prefix for file name.
func WithFilePrefixStringer(s fmt.Stringer) WrapOption {
	return wrapOptionFunc(func(o *Option) {
		o.prefix = s
	})
}
