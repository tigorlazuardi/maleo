package maleos3

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Option interface {
	apply(*S3)
}

type OptionFunc func(*S3)

func (f OptionFunc) apply(s *S3) {
	f(s)
}

// WithClient sets the S3 client to use.
func WithClient(client *s3.Client) Option {
	return OptionFunc(func(s *S3) {
		s.Client = client
	})
}

// WithRegion sets the region to use.
func WithRegion(region string) Option {
	return OptionFunc(func(s *S3) {
		s.region = region
	})
}

// WithBucket sets the bucket name to use.
func WithBucket(name string) Option {
	return OptionFunc(func(s *S3) {
		s.bucket = name
	})
}

// WithSecure sets to use https or not.
func WithSecure(secure bool) Option {
	return OptionFunc(func(s *S3) {
		s.secure = secure
	})
}

// WithURLBuilder sets the url builder to use.
func WithURLBuilder(builder URLBuilder) Option {
	return OptionFunc(func(s *S3) {
		s.urlBuilder = builder
	})
}

type StringerFunc func() string

func (f StringerFunc) String() string {
	return f()
}

// WithFilenamePretext sets the pretext to use for the filename to be uploaded to the bucket.
func WithFilenamePretext(pretext string) Option {
	return OptionFunc(func(s *S3) {
		s.pretext = StringerFunc(func() string {
			return pretext
		})
	})
}

// WithFilenameDynamicPretext sets the pretext to use for the filename to be uploaded to the bucket.
//
// example of creating yearly folder in bucket:
//
//	WithFilenameDynamicPretext(func() string { return time.Now().Format("2006") + "/" })
func WithFilenameDynamicPretext(pretext func() string) Option {
	return OptionFunc(func(s *S3) {
		s.pretext = StringerFunc(pretext)
	})
}
