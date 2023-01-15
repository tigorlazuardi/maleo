package maleo

import (
	"context"
	"sync"
)

type mockLogger struct {
	called bool
}

func newMockLogger() *mockLogger {
	return &mockLogger{}
}

func (m *mockLogger) Log(ctx context.Context, entry Entry) {
	m.called = true
}

func (m *mockLogger) LogError(ctx context.Context, err Error) {
	m.called = true
}

func newMockMessenger(count int) *mockMessenger {
	wg := &sync.WaitGroup{}
	wg.Add(count)
	return &mockMessenger{
		wg: wg,
	}
}

func newMockMessengerWithName(count int, name string) *mockMessenger {
	wg := &sync.WaitGroup{}
	wg.Add(count)
	return &mockMessenger{
		wg:   wg,
		name: name,
	}
}

type mockMessenger struct {
	called bool
	wg     *sync.WaitGroup
	name   string
}

func (m *mockMessenger) String() string {
	return m.name
}

func (m *mockMessenger) Name() string {
	if m.name == "" {
		return "mock"
	}
	return m.name
}

func (m *mockMessenger) SendMessage(ctx context.Context, msg MessageContext) {
	m.called = true
	m.wg.Done()
}

func (m *mockMessenger) Wait(ctx context.Context) error {
	m.wg.Wait()
	return nil
}
