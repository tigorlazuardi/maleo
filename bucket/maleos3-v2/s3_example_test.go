package maleos3_test

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/tigorlazuardi/maleo/bucket"
	"github.com/tigorlazuardi/maleo/bucket/maleos3-v2"
)

// --8<-- [start:integrated]

func ExampleNewS3Bucket() {
	bkt, err := maleos3.NewS3Bucket("my-bucket.s3.us-east-1.amazonaws.com")
	if err != nil {
		return
	}
	f := strings.NewReader("hello world")
	file := bucket.NewFile(f, "text/plain; charset=utf-8")
	for _, result := range bkt.Upload(context.Background(), []bucket.File{file}) {
		if result.Error != nil {
			// handle error
			return
		}
	}
}

// --8<-- [end:integrated]
// --8<-- [start:static]

func ExampleNewS3Bucket_static() {
	client, err := maleos3.SimpleStaticClient(maleos3.SimpleStaticParams{
		AccessKeyID:     "access_key",
		SecretAccessKey: "secret_key",
		SessionToken:    "", // Optional. Leave empty if not needed.

		// Endpoint is needed to be set for non-AWS S3, otherwise just leave empty.
		// It is needed as interoperability with other S3 compatible services.
		Endpoint: "",
		// Same reason as above
		Region: "",
		Secure: true, // Set false if you want to use HTTP instead of HTTPS.
	})
	if err != nil {
		return
	}
	bkt, err := maleos3.NewS3Bucket("my-bucket.s3.us-east-1.amazonaws.com", maleos3.WithClient(client))
	if err != nil {
		return
	}
	f := strings.NewReader("hello world")
	file := bucket.NewFile(f, "text/plain; charset=utf-8")
	for _, result := range bkt.Upload(context.Background(), []bucket.File{file}) {
		if result.Error != nil {
			// handle error
			return
		}
	}
}

// --8<-- [end:static]
// --8<-- [start:custom]

func ExampleNewS3Bucket_custom() {
	client := s3.NewFromConfig(aws.Config{
		// set your custom config here
	})
	bkt, err := maleos3.NewS3Bucket("bucket.s3.us-east-1.amazonaws.com", maleos3.WithClient(client))
	if err != nil {
		return
	}
	f := strings.NewReader("hello world")
	file := bucket.NewFile(f, "text/plain; charset=utf-8")
	for _, result := range bkt.Upload(context.Background(), []bucket.File{file}) {
		if result.Error != nil {
			// handle error
			return
		}
	}
}

// --8<-- [end:custom]

// --8<-- [start:nonaws]

func ExampleNewS3Bucket_nonAws() {
	client, err := maleos3.SimpleStaticClient(maleos3.SimpleStaticParams{
		AccessKeyID:     "access_key",
		SecretAccessKey: "secret_key",
		Endpoint:        "minio:9000",
		Region:          "us-east-1",
	})
	if err != nil {
		return
	}
	bkt, err := maleos3.NewS3Bucket("minio:9000", maleos3.WithClient(client))
	if err != nil {
		return
	}
	f := strings.NewReader("hello world")
	file := bucket.NewFile(f, "text/plain; charset=utf-8")
	for _, result := range bkt.Upload(context.Background(), []bucket.File{file}) {
		if result.Error != nil {
			// handle error
			return
		}
	}
}

// --8<-- [end:nonaws]
