package maleo

type Engine interface {
	ErrorConstructor
	EntryConstructor
	EntryMessageContextConstructor
	ErrorMessageContextConstructor
}

func NewEngine(opts ...EngineOption) Engine {
	def := &engine{
		ErrorConstructor:               ErrorConstructorFunc(defaultErrorGenerator),
		EntryConstructor:               EntryConstructorFunc(defaultEntryConstructor),
		EntryMessageContextConstructor: MessageContextConstructorFunc(defaultMessageContextConstructor),
		ErrorMessageContextConstructor: ErrorMessageConstructorFunc(defaultErrorMessageContextConstructor),
	}
	for _, opt := range opts {
		opt.apply(def)
	}
	return def
}

type engine struct {
	ErrorConstructor
	EntryConstructor
	EntryMessageContextConstructor
	ErrorMessageContextConstructor
}
