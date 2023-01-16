package maleos3

import (
	"context"
	"fmt"
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
	checkEnvs(t, "AWS_ENDPOINT", "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY")
	bkt, err := NewS3Bucket(os.Getenv("AWS_ENDPOINT"))
	if err != nil {
		return
	}
	f := strings.NewReader("hello world")
	file := bucket.NewFile(f, "text/plain; charset=utf-8", bucket.WithFilename("test.txt"))
	results := bkt.Upload(context.Background(), []bucket.File{file})
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	for _, result := range results {
		if result.Error != nil {
			t.Fatalf("unexpected error: %v", result.Error)
		}
		url := fmt.Sprintf("https://%s/%s", os.Getenv("AWS_ENDPOINT"), result.File.Filename())
		if result.URL != url {
			t.Fatalf("unexpected url: %s", result.URL)
		}
		if t.Failed() {
			t.Logf("failed to upload file: %s to %s", result.File.Filename(), result.URL)
		}
	}
}
