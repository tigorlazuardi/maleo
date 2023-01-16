package maleominio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/tigorlazuardi/maleo/bucket"
	"time"
)

type Minio struct {
	*minio.Client
	bucket string
	opts   *Option
}

func (m Minio) Upload(ctx context.Context, files []bucket.File) []bucket.UploadResult {
	results := make([]bucket.UploadResult, 0, len(files))
	for _, f := range files {
		func(f bucket.File) {
			defer func() { _ = f.Close() }()
			filename := m.opts.prefix.String() + f.Filename()
			objectURL := m.Client.EndpointURL().String() + "/" + m.bucket + "/" + filename
			_, err := m.PutObject(ctx, m.bucket, filename, f, int64(f.Size()), m.opts.putObjectOptions(ctx, f))
			results = append(results, bucket.UploadResult{
				URL:   objectURL,
				File:  f,
				Error: err,
			})
		}(f)
	}
	return results
}

// Wrap wraps minio client to tower bucket implementation.
//
// Client credentials must have permission to write access to target bucket.
//
// If client has permission to list and create bucket, it will be used to check if bucket is not exist on initialization,
// and create the bucket if it is not exist. If client does not have such permissions, it will be silently ignored.
//
// You may operate how bucket creation and put object operation is executed by providing options.
func Wrap(client *minio.Client, bucketName string, options ...WrapOption) bucket.Bucket {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opts := &Option{
		prefix: StringerFunc(func() string {
			return ""
		}),
		putObjectOptions: func(ctx context.Context, file bucket.File) minio.PutObjectOptions {
			return minio.PutObjectOptions{
				ContentType: file.ContentType(),
			}
		},
	}
	for _, o := range options {
		o.apply(opts)
	}
	if buckets, err := client.ListBuckets(ctx); err == nil {
		for _, b := range buckets {
			if b.Name == bucketName {
				return Minio{client, bucketName, opts}
			}
		}
		_ = client.MakeBucket(ctx, bucketName, opts.makeBucketOption)
	}
	return &Minio{client, bucketName, opts}
}
