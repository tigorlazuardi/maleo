package maleo

type option struct{}

var Option option

func (option) Init() InitOptionBuilder {
	return InitOptionBuilder{}
}

func (option) Engine() EngineOptionBuilder {
	return EngineOptionBuilder{}
}

func (option) Message() MessageOptionBuilder {
	return MessageOptionBuilder{}
}
