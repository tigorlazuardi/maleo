package maleogomemcache

import (
	"context"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/tigorlazuardi/maleo/locker"
)

type Locker struct {
	client *memcache.Client
}

func (l Locker) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	if len(key) > 250 {
		key = key[:250]
	}
	item := &memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: int32(ttl.Seconds()),
	}
	err := l.client.Set(item)
	if err != nil {
		return fmt.Errorf("unable to set value to key '%s': %w", key, err)
	}
	return nil
}

func (l Locker) Get(_ context.Context, key string) ([]byte, error) {
	if len(key) > 250 {
		key = key[:250]
	}
	item, err := l.client.Get(key)
	if err != nil {
		return nil, fmt.Errorf("unable to get value from key '%s': %w", key, err)
	}
	return item.Value, nil
}

func (l Locker) Delete(_ context.Context, key string) {
	if len(key) > 250 {
		key = key[:250]
	}
	_ = l.client.Delete(key)
}

func (l Locker) Exist(_ context.Context, key string) bool {
	if len(key) > 250 {
		key = key[:250]
	}
	_, err := l.client.Get(key)
	return err == nil
}

func (l Locker) Separator() string {
	return "::"
}

func New(client *memcache.Client) locker.Locker {
	return Locker{client: client}
}
