package maleogoredis

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/go-redis/redis/v9"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func createClient() (*redis.Client, func(), error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(err)
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "7-alpine",
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
	target := fmt.Sprintf("%s:%s", host, resource.GetPort("6379/tcp"))
	var client *redis.Client
	if err := pool.Retry(func() error {
		client = redis.NewClient(&redis.Options{
			Addr: target,
		})
		return client.Ping(context.Background()).Err()
	}); err != nil {
		_ = pool.Purge(resource)
		return nil, nil, fmt.Errorf("could not connect to docker target '%s': %w", target, err)
	}
	cleanup := func() {
		_ = pool.Purge(resource)
	}
	return client, cleanup, nil
}

func TestGoRedisV9(t *testing.T) {
	client, cleanup, err := createClient()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	key := strings.Repeat("foo", 250)
	cache := New(client)
	ctx := context.Background()
	err = cache.Set(ctx, key, []byte("bar"), 0)
	if err != nil {
		t.Fatal(err)
	}
	val, err := cache.Get(ctx, key)
	if err != nil {
		t.Fatal(err)
	}
	if string(val) != "bar" {
		t.Fatalf("expected 'bar', got '%s'", string(val))
	}
	cache.Delete(ctx, key)
	if cache.Exist(ctx, key) {
		t.Fatal("expected key to be deleted")
	}
	if cache.Separator() != ":" {
		t.Fatalf("expected separator to be ':', got '%s'", cache.Separator())
	}

	_, err = cache.Get(ctx, "not-exist")
	if err == nil {
		t.Fatal("expected error when key is not exist")
	}
}
