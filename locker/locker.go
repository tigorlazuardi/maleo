package locker

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"
)

var ErrNil = errors.New("value does not exist")

// --8<-- [start:locker]

type Locker interface {
	// Set the key and value.
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	// Get the Value by Key.
	Get(ctx context.Context, key string) ([]byte, error)
	// Delete value by key.
	Delete(ctx context.Context, key string)
	// Exist Checks if key exist in cache.
	Exist(ctx context.Context, key string) bool
	// Separator Returns Accepted separator value for the Lock implementor.
	Separator() string
}

// --8<-- [end:locker]

type lockValue struct {
	value []byte
	time  time.Time
}

var _ Locker = (*LocalLock)(nil)

type LocalLock struct {
	mu            *sync.RWMutex
	state         map[string]*lockValue
	length        int
	lastRebalance time.Time
}

// NewLocalLock creates a RAM based cache.
//
// This cache is not persistent and will be lost on application restart.
//
// This cache does not support distributed caching mechanism, and is not safe for multiple application
// that uses the same key for handling rate limits
//
// (e.g. Discord enforced 1 second limit between messages with the same token, if the token is shared between services,
// and you use this local cache for all your machines, your token may be banned by Discord for over limit since the services
// cannot know the state of current rate limit, and thus just assume everything is safe to be sent).
//
// Use this cache for tests or when you know that you will not have multiple application instances.
func NewLocalLock() *LocalLock {
	return &LocalLock{
		mu:            &sync.RWMutex{},
		state:         make(map[string]*lockValue),
		length:        0,
		lastRebalance: time.Now(),
	}
}

// Set Sets the Cache key and value.
func (m *LocalLock) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	if ttl < 1 {
		ttl = math.MaxInt64
	}
	m.mu.Lock()
	cache := m.state[key]
	if cache == nil {
		m.length += 1
		cache = &lockValue{}
	}
	cache.value = value
	cache.time = time.Now().Add(ttl)
	m.state[key] = cache
	m.mu.Unlock()
	m.checkGC()
	return nil
}

// Get the Value by Key. Returns maleo.ErrNilCache if not found or ttl has passed.
func (m *LocalLock) Get(ctx context.Context, key string) ([]byte, error) {
	m.mu.RLock()
	cache, ok := m.state[key]
	if !ok {
		m.mu.RUnlock()
		return nil, ErrNil
	}
	m.mu.RUnlock()
	now := time.Now()
	if now.After(cache.time) {
		m.Delete(ctx, key)
		return nil, ErrNil
	}
	return cache.value, nil
}

// Exist Checks if Key exist in cache.
func (m *LocalLock) Exist(_ context.Context, key string) bool {
	m.mu.RLock()
	_, ok := m.state[key]
	m.mu.RUnlock()
	return ok
}

// Delete key from cache.
func (m *LocalLock) Delete(_ context.Context, key string) {
	m.mu.Lock()
	delete(m.state, key)
	if m.length > 0 {
		m.length -= 1
	}
	m.mu.Unlock()
}

func (m *LocalLock) checkGC() {
	now := time.Now()
	if now.After(m.lastRebalance.Add(time.Minute*5)) && m.length > 1000 {
		go func() {
			m.mu.Lock()
			n := make(map[string]*lockValue, len(m.state))
			for k, v := range m.state {
				n[k] = v
			}
			m.state = n
			m.lastRebalance = time.Now()
			m.length = len(n)
			m.mu.Unlock()
		}()
	}
}

func (m *LocalLock) Separator() string {
	return "::"
}
