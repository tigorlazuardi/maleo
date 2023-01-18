package maleominio_test

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/tigorlazuardi/maleo/bucket"
	"github.com/tigorlazuardi/maleo/bucket/maleominio-v7"
)

// --8<-- [start:example]

func ExampleWrap() {
	client, err := minio.New("play.min.io", &minio.Options{
		Creds:  credentials.NewStaticV4("access_key", "secret_key", ""),
		Secure: true,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	bkt := maleominio.Wrap(client, "my-bucket")
	f := strings.NewReader("hello world")
	file := bucket.NewFile(f, "text/plain; charset=utf-8")
	for _, result := range bkt.Upload(context.Background(), []bucket.File{file}) {
		if result.Error != nil {
			fmt.Println(result.Error.Error())
			return
		}
	}
}

// --8<-- [end:example]

// --8<-- [start:aws]

func ExampleWrap_aws() {
	// Note the region. It's mandatory.
	client, err := minio.New("s3.us-east-1.amazonaws.com", &minio.Options{
		Creds:  credentials.NewStaticV4("access_key", "secret_key", ""),
		Secure: true,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	bkt := maleominio.Wrap(client, "my-bucket",
		// add optional "yyyy-mm-dd" folder to all uploaded files.
		maleominio.WithFilePrefixStringer(
			maleominio.StringerFunc(func() string {
				return time.Now().Format("2006-01-02/")
			}),
		),
	)
	f := strings.NewReader("hello world")
	file := bucket.NewFile(f, "text/plain; charset=utf-8")
	for _, result := range bkt.Upload(context.Background(), []bucket.File{file}) {
		if result.Error != nil {
			fmt.Println(result.Error.Error())
			return
		}
	}
}

// --8<-- [end:aws]
