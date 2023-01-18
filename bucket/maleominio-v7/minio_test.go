package maleominio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/tigorlazuardi/maleo/bucket"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func createClient() (*minio.Client, func(), error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(err)
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "quay.io/minio/minio",
		Tag:        "RELEASE.2023-01-02T09-40-09Z",
		Cmd:        []string{"server", "/data"},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		return nil, nil, err
	}
	host := "172.17.0.1" // docker0 interface
	if os.Getenv("DOCKER_TEST_HOST") != "" {
		host = os.Getenv("DOCKER_TEST_HOST")
	}
	target := fmt.Sprintf("%s:%s", host, resource.GetPort("9000/tcp"))
	var client *minio.Client
	if err := pool.Retry(func() error {
		var err error
		client, err = minio.New(target, &minio.Options{
			Creds: credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		})
		if err != nil {
			fmt.Println(err.Error())
		}
		return err
	}); err != nil {
		_ = pool.Purge(resource)
		return nil, nil, fmt.Errorf("could not connect to docker target '%s': %w", target, err)
	}
	cleanup := func() {
		_ = pool.Purge(resource)
	}
	return client, cleanup, nil
}

const policy = `
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": [
          "*"
        ]
      },
      "Action": [
        "s3:GetBucketLocation",
        "s3:ListBucket"
      ],
      "Resource": [
        "arn:aws:s3:::test"
      ]
    },
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": [
          "*"
        ]
      },
      "Action": [
        "s3:GetObject"
      ],
      "Resource": [
        "arn:aws:s3:::test/*"
      ]
    }
  ]
}`

func TestMinio_Upload(t *testing.T) {
	const test = "test"
	client, clean, err := createClient()
	if err != nil {
		t.Fatal(err)
	}
	defer clean()
	wc := Wrap(client, test,
		WithMakeBucketOption(minio.MakeBucketOptions{}),
		WithPutObjectOption(func(ctx context.Context, file bucket.File) minio.PutObjectOptions {
			return minio.PutObjectOptions{
				ContentType: file.ContentType(),
			}
		}),
		WithFilePrefixStringer(StringerFunc(func() string {
			return "prefix/"
		})))
	f := bucket.NewFile(strings.NewReader(test), "text/plain; charset=utf-8",
		bucket.WithFilename("test.txt"),
		bucket.WithFilesize(len(test)),
	)
	err = client.SetBucketPolicy(context.Background(), test, policy)
	if err != nil {
		t.Fatalf("could not set bucket policy: %v", err)
	}
	results := wc.Upload(context.Background(), []bucket.File{f})
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	for _, result := range results {
		func(result bucket.UploadResult) {
			if result.Error != nil {
				t.Fatalf("unexpected error: %v", result.Error)
			}
			resp, err := http.Get(result.URL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			defer func(Body io.ReadCloser) {
				if err := Body.Close(); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}(resp.Body)
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("unexpected status code: %d", resp.StatusCode)
			}
			if resp.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
				t.Fatalf("unexpected content type: %s", resp.Header.Get("Content-Type"))
			}
			content, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("could not read object: %v", err)
			}
			if string(content) != "test" {
				t.Fatalf("unexpected content: %s", content)
			}
		}(result)
	}
	wc = Wrap(client, test, WithFilePrefix("test/"))
	f = bucket.NewFile(strings.NewReader(test), "text/plain; charset=utf-8",
		bucket.WithFilename("test.txt"),
		bucket.WithFilesize(len(test)),
	)
	results = wc.Upload(context.Background(), []bucket.File{f})
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	for _, result := range results {
		func(result bucket.UploadResult) {
			if result.Error != nil {
				t.Fatalf("unexpected error: %v", result.Error)
			}
			obj, err := client.GetObject(context.Background(), "test", "test/"+result.File.Filename(), minio.GetObjectOptions{})
			if err != nil {
				t.Fatalf("could not get object: %v", err)
			}
			defer func(obj *minio.Object) {
				err := obj.Close()
				if err != nil {
					t.Fatalf("could not close object: %v", err)
				}
			}(obj)
		}(result)
	}

	wc = Wrap(client, test)
	f = bucket.NewFile(strings.NewReader(test), "text/plain; charset=utf-8",
		bucket.WithFilename("test.txt"),
		bucket.WithFilesize(len(test)),
	)
	results = wc.Upload(context.Background(), []bucket.File{f})
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	for _, result := range results {
		func(result bucket.UploadResult) {
			if result.Error != nil {
				t.Fatalf("unexpected error: %v", result.Error)
			}
			obj, err := client.GetObject(context.Background(), "test", result.File.Filename(), minio.GetObjectOptions{})
			if err != nil {
				t.Fatalf("could not get object: %v", err)
			}
			defer func(obj *minio.Object) {
				err := obj.Close()
				if err != nil {
					t.Fatalf("could not close object: %v", err)
				}
			}(obj)
		}(result)
	}
}
