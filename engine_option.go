package maleo

type EngineOption interface {
	apply(*engine)
}

type (
	EngineOptionFunc    func(*engine)
	EngineOptionBuilder []EngineOption
)

func (e EngineOptionFunc) apply(m *engine) {
	e(m)
}

func (e EngineOptionBuilder) apply(m *engine) {
	for _, opt := range e {
		opt.apply(m)
	}
}

// ErrorConstructor sets the error constructor for the engine.
func (e EngineOptionBuilder) ErrorConstructor(ec ErrorConstructor) EngineOptionBuilder {
	return append(e, EngineOptionFunc(func(m *engine) {
		m.ErrorConstructor = ec
	}))
}

// EntryConstructor sets the entry constructor for the engine.
func (e EngineOptionBuilder) EntryConstructor(ec EntryConstructor) EngineOptionBuilder {
	return append(e, EngineOptionFunc(func(m *engine) {
		m.EntryConstructor = ec
	}))
}

// EntryMessageContextConstructor sets the message context constructor for the engine.
func (e EngineOptionBuilder) EntryMessageContextConstructor(mc EntryMessageContextConstructor) EngineOptionBuilder {
	return append(e, EngineOptionFunc(func(m *engine) {
		m.EntryMessageContextConstructor = mc
	}))
}

func (e EngineOptionBuilder) ErrorMessageContextConstructor(mc ErrorMessageContextConstructor) EngineOptionBuilder {
	return append(e, EngineOptionFunc(func(m *engine) {
		m.ErrorMessageContextConstructor = mc
	}))
}
