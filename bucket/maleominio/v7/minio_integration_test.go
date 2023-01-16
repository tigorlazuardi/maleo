package maleominio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/tigorlazuardi/maleo/bucket"
	"github.com/tigorlazuardi/maleo/loader"
	"os"
	"strings"
	"testing"
)

func checkEnvs(t *testing.T, envs ...string) {
	for _, env := range envs {
		if os.Getenv(env) == "" {
			t.Skip(env + " env not found. skipping integration test")
		}
	}
}

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	loader.LoadEnv()
	checkEnvs(t, "MINIO_TEST_ENDPOINT", "MINIO_TEST_ACCESS_KEY", "MINIO_TEST_SECRET_KEY", "MINIO_TEST_BUCKET")
	client, err := minio.New(os.Getenv("MINIO_TEST_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_TEST_ACCESS_KEY"), os.Getenv("MINIO_TEST_SECRET_KEY"), ""),
		Secure: true,
	})
	if err != nil {
		t.Fatalf("could not create minio client: %v", err)
	}
	m := Wrap(client, os.Getenv("MINIO_TEST_BUCKET"))
	file := bucket.NewFile(strings.NewReader("hello world"), "text/plain; charset=utf-8", bucket.WithFilename("test.txt"))
	results := m.Upload(context.Background(), []bucket.File{file})
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	for _, result := range results {
		if result.Error != nil {
			t.Fatalf("unexpected error: %v", result.Error)
		}
	}
}
