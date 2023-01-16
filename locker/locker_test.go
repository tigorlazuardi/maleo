package locker

import (
	"strconv"
	"testing"
	"time"
)

func TestLocalLock(t *testing.T) {
	cache := NewLocalLock()
	err := cache.Set(nil, "test", []byte("test"), 0)
	if err != nil {
		t.Fatal(err)
	}
	if !cache.Exist(nil, "test") {
		t.Fatal("key should exist")
	}
	value, err := cache.Get(nil, "test")
	if err != nil {
		t.Fatal(err)
	}
	if string(value) != "test" {
		t.Fatal("value should be test")
	}
	cache.Delete(nil, "test")
	if cache.Exist(nil, "test") {
		t.Fatal("key should not exist")
	}

	_, err = cache.Get(nil, "test")
	if err != Nil {
		t.Fatal("should return ErrNilCache")
	}
	err = cache.Set(nil, "test", []byte("test"), time.Nanosecond)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Nanosecond)
	_, err = cache.Get(nil, "test")
	if err != Nil {
		t.Fatal("should return ErrNilCache")
	}

	for i := 0; i < 1001; i++ {
		_ = cache.Set(nil, strconv.Itoa(i), []byte("test"), 0)
	}
	if cache.Separator() != "::" {
		t.Fatal("separator should be ::")
	}
	cache.lastRebalance = time.Now().Add(-time.Hour)
	cache.checkGC()
	time.Sleep(time.Millisecond)
}
