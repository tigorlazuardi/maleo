package maleo

var global *Maleo

func init() {
	global = NewMaleo(Service{}, Option.Init().CallerDepth(3).Name("maleo-global"))
}

// Global returns the global Maleo instance.
func Global() *Maleo {
	return global
}

// SetGlobal sets the global Maleo instance.
func SetGlobal(m *Maleo) {
	global = m
}
