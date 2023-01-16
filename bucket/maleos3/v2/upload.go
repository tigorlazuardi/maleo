package maleos3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/tigorlazuardi/maleo/bucket"
	"sync"
)

func (s *S3) Upload(ctx context.Context, files []bucket.File) []bucket.UploadResult {
	results := make([]bucket.UploadResult, len(files))
	wg := &sync.WaitGroup{}
	wg.Add(len(files))
	for i, file := range files {
		go func(i int, file bucket.File) {
			defer wg.Done()
			defer func() { _ = file.Close() }()
			filename := file.Filename()
			pretext := s.pretext.String()
			if pretext != "" {
				filename = pretext + filename
			}
			var opts []func(uploader *manager.Uploader)
			if s.resolver != nil {
				opts = append(opts, manager.WithUploaderRequestOptions(s3.WithEndpointResolver(s.resolver)))
			}

			_, err := s.uploader.Upload(ctx, &s3.PutObjectInput{
				Bucket:      aws.String(s.bucket),
				Key:         aws.String(filename),
				Body:        file,
				ContentType: aws.String(file.ContentType()),
			}, opts...)

			proto := "https"
			if !s.secure {
				proto = "http"
			}
			obj := NewObject(proto, s.endpoint, s.region, s.bucket, filename)
			results[i] = bucket.UploadResult{
				URL:   s.urlBuilder.BuildURL(obj),
				File:  file,
				Error: err,
			}
		}(i, file)
	}
	wg.Wait()
	return results
}
