package maleo

import (
	"context"
	"reflect"
	"testing"
)

func TestMessengers_Filter(t *testing.T) {
	type args struct {
		f FilterMessengersFunc
	}
	tests := []struct {
		name string
		m    Messengers
		args args
		want Messengers
	}{
		{
			name: "empty",
			m:    Messengers{},
			args: args{
				f: func(m Messenger) bool {
					return true
				},
			},
			want: Messengers{},
		},
		{
			name: "single",
			m:    Messengers{newMockMessenger(0)},
			args: args{
				f: func(messenger Messenger) bool {
					return messenger.Name() == "mock"
				},
			},
			want: Messengers{newMockMessenger(0)},
		},
		{
			name: "single - filter",
			m:    Messengers{newMockMessenger(0)},
			args: args{
				f: func(messenger Messenger) bool {
					return messenger.Name() == "mock2"
				},
			},
			want: Messengers{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Filter(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessengers_SendMessage(t *testing.T) {
	type args struct {
		ctx context.Context
		msg MessageContext
	}
	tests := []struct {
		name string
		m    Messengers
		args args
	}{
		{
			name: "called",
			m:    Messengers{newMockMessenger(1), newMockMessenger(1)},
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.SendMessage(tt.args.ctx, tt.args.msg)
			if err := tt.m.Wait(context.Background()); err != nil {
				t.Fatalf("Wait() error = %v", err)
			}
			for _, messenger := range tt.m {
				if !messenger.(*mockMessenger).called {
					t.Errorf("SendMessage() = %v, want %v", messenger.(*mockMessenger).called, 1)
				}
			}
		})
	}
}
