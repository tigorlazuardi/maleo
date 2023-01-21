package maleo

import "context"

var global *Maleo

func init() {
	global = New(Service{}, Option.Init().CallerDepth(3).Name("maleo-global"))
}

// Global returns the global Maleo instance.
func Global() *Maleo {
	return global
}

// SetGlobal sets the global Maleo instance.
func SetGlobal(m *Maleo) {
	global = m
}

// Wait waits for ongoing messages for Messenger to finish.
func Wait(ctx context.Context) error {
	return global.Wait(ctx)
}
