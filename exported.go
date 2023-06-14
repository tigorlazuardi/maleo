package maleo

import "context"

var global *Maleo

func init() {
	global = New(Service{}, Option.Init().CallerDepth(3).Name("maleo-global"))
	global.isGlobal = true
}

// Global returns the global Maleo instance.
//
// To get a copy of the global Maleo instance, use maleo.Global().Clone()
//
// Do not attach Global instance to struct members or somewhere long lived, call Clone() first and use the cloned instance instead.
//
// Cloned Maleo instance will ignore Maleo instances attached to context.
func Global() *Maleo {
	return global
}

// SetGlobal sets the global Maleo instance.
//
// Maleo instance will be set to receive overrides to use Maleo instances attached to context.
func SetGlobal(m *Maleo) {
	m.isGlobal = true
	global = m
}

// Wait waits for ongoing messages for Messenger to finish.
func Wait(ctx context.Context) error {
	return global.Wait(ctx)
}
