package maleo

import (
	"context"
	"errors"
	"testing"
)

func Test_multierror_Error(t *testing.T) {
	tests := []struct {
		name string
		m    multierror
		want string
	}{
		{
			name: "empty",
			m:    multierror{},
			want: "",
		},
		{
			name: "single",
			m:    multierror{errors.New("test")},
			want: "1. test",
		},
		{
			name: "multiple",
			m:    multierror{errors.New("test"), errors.New("test2")},
			want: "1. test; 2. test2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMaleo(t *testing.T) {
	mal := NewMaleo(Service{}, Option.Init().
		Logger(NoopLogger{}).
		Name("test").
		CallerDepth(1).
		DefaultMessageOptions(Option.Message().ForceSend(true)).
		Messengers(newMockMessengerWithName(2, "test1")),
	)
	mal.SetLogger(NoopLogger{})
	mal.SetEngine(NewEngine())
	if mal.Name() != "test" {
		t.Errorf("Name() = %v, want %v", mal.Name(), "test")
	}
	if mal.callerDepth != 1 {
		t.Errorf("callerDepth = %v, want %v", mal.callerDepth, 1)
	}
	if mal.logger == nil {
		t.Errorf("logger = %v, want %v", mal.logger, NoopLogger{})
	}
	if mal.engine == nil {
		t.Errorf("engine = %v, want %v", mal.engine, NewEngine())
	}
	mal.RegisterBenched(newMockMessengerWithName(0, "test3"))
	mal.LogError(context.Background(), mal.BailFreeze("foo"))
	mal.Log(context.Background(), mal.NewEntry("foo %s", "bar").Freeze())
	mal.NotifyError(context.Background(), mal.Bail("foo %s", "bar").Freeze(), Option.Message().ForceSend(true))
	mal.Notify(context.Background(), mal.NewEntry("foo").Freeze(), Option.Message().ForceSend(true))
	if err := mal.Wait(context.Background()); err != nil {
		t.Errorf("Wait() = %v, want %v", err, nil)
	}
	for _, v := range mal.defaultParams.Messengers {
		v := v.(*mockMessenger)
		if !v.called {
			t.Errorf("messenger %v is not called", v.name)
		}
	}
	if mal.Service() != (Service{}) {
		t.Errorf("Service() = %v, want %v", mal.Service(), Service{})
	}
}
