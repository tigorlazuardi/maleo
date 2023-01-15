package maleo

type InitOption interface {
	apply(*Maleo)
}

type (
	InitOptionFunc    func(*Maleo)
	InitOptionBuilder []InitOption
)

func (i InitOptionFunc) apply(m *Maleo) {
	i(m)
}

func (i InitOptionBuilder) apply(m *Maleo) {
	for _, opt := range i {
		opt.apply(m)
	}
}

// Logger sets the logger for Maleo.
func (i InitOptionBuilder) Logger(l Logger) InitOptionBuilder {
	return append(i, InitOptionFunc(func(m *Maleo) {
		m.logger = l
	}))
}

// DefaultMessageOptions overrides the default options for messages.
func (i InitOptionBuilder) DefaultMessageOptions(opts ...MessageOption) InitOptionBuilder {
	return append(i, InitOptionFunc(func(m *Maleo) {
		for _, opt := range opts {
			opt.Apply(m.defaultParams)
		}
	}))
}

// Messengers sets the default messengers for Maleo.
func (i InitOptionBuilder) Messengers(messengers ...Messenger) InitOptionBuilder {
	return append(i, InitOptionFunc(func(m *Maleo) {
		m.defaultParams.Messengers = messengers
	}))
}

// Name sets the name of the Maleo instance. This is used for managing the messengers.
//
// Only useful if you set this Maleo instance as a Messenger for other Maleo instances.
func (i InitOptionBuilder) Name(name string) InitOptionBuilder {
	return append(i, InitOptionFunc(func(m *Maleo) {
		m.name = name
	}))
}

func (i InitOptionBuilder) CallerDepth(depth int) InitOptionBuilder {
	return append(i, InitOptionFunc(func(m *Maleo) {
		m.callerDepth = depth
	}))
}
