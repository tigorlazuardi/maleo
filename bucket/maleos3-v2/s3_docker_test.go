package maleos3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/tigorlazuardi/maleo/bucket"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func createClient() (string, func(), error) {
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
		return "", nil, err
	}
	host := "172.17.0.1" // docker0 interface
	if os.Getenv("DOCKER_TEST_HOST") != "" {
		host = os.Getenv("DOCKER_TEST_HOST")
	}
	var target string
	if err := pool.Retry(func() error {
		target = fmt.Sprintf("%s:%s", host, resource.GetPort("9000/tcp"))
		resp, err := http.Get("http://" + target)
		if err != nil {
			return err
		}
		_ = resp.Body.Close()
		return nil
	}); err != nil {
		_ = pool.Purge(resource)
		return "", nil, fmt.Errorf("could not connect to docker target '%s': %w", host, err)
	}
	cleanup := func() {
		_ = pool.Purge(resource)
	}
	return target, cleanup, nil
}

func createPublicPolicy(bucketName string) string {
	return fmt.Sprintf(`{
	  "Version":"2012-10-17",
	  "Statement":[
		{
		  "Effect":"Allow",
		  "Principal":{
			"AWS":["*"]
		  },
		  "Action":[
			"s3:GetBucketLocation",
			"s3:ListBucket"
		  ],
		  "Resource":[
			"arn:aws:s3:::%s"
		  ]
		},
		{
		  "Effect":"Allow",
		  "Principal":{
			"AWS":["*"]
		  },
		  "Action":["s3:GetObject"],
		  "Resource":["arn:aws:s3:::%s/*"]
		}
	  ]
	}`, bucketName, bucketName)
}

func TestS3(t *testing.T) {
	host, cleanup, err := createClient()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	t.Run("Upload", s3UploadTest(host))
}

func s3UploadTest(host string) func(t *testing.T) {
	return func(t *testing.T) {
		client, err := SimpleStaticClient(SimpleStaticParams{
			AccessKeyID:     "minioadmin",
			SecretAccessKey: "minioadmin",
			SessionToken:    "",
			Endpoint:        host,
			Region:          "us-east-1",
			Secure:          false,
		})
		if err != nil {
			t.Fatalf("failed to create config: %v", err)
		}
		now := time.Now()
		date := now.Format("2006-01-02")
		bkt, err := NewS3Bucket(host,
			WithBucket("my-bucket"),
			WithSecure(false),
			WithRegion("us-east-1"),
			WithFilenamePretext("foo"),
			WithFilenameDynamicPretext(func() string {
				return date + "/"
			}),
			WithClient(client),
		)
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.PutBucketPolicy(context.Background(), &s3.PutBucketPolicyInput{
			Bucket: aws.String("my-bucket"),
			Policy: aws.String(createPublicPolicy("my-bucket")),
		})
		if err != nil {
			t.Fatalf("failed to update bucket ACL: %v", err)
		}
		file := bucket.NewFile(strings.NewReader("hello world"), "text/plain; charset=utf-8", bucket.WithFilename("hello.txt"))
		results := bkt.Upload(context.Background(), []bucket.File{file})
		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		result := results[0]
		if result.Error != nil {
			t.Fatal(result.Error)
		}
		targetName := fmt.Sprintf("http://%s/my-bucket/%s/hello.txt", host, date)
		if result.URL != targetName {
			t.Fatalf("expected %s, got %s", targetName, result.URL)
		}
		resp, err := http.Get(result.URL)
		if err != nil {
			t.Fatalf("failed to get file: %v", err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Fatalf("failed to close body: %v", err)
			}
		}(resp.Body)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}
		if string(b) != "hello world" {
			t.Fatalf("expected 'hello world', got %s", string(b))
		}
	}
}
