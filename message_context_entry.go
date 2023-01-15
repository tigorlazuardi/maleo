package maleo

import "time"

type EntryMessageContextConstructor interface {
	BuildEntryMessageContext(entry Entry, param *MessageParameters) MessageContext
}

type MessageContextConstructorFunc func(entry Entry, param *MessageParameters) MessageContext

func (m MessageContextConstructorFunc) BuildEntryMessageContext(entry Entry, param *MessageParameters) MessageContext {
	return m(entry, param)
}

func defaultMessageContextConstructor(entry Entry, param *MessageParameters) MessageContext {
	return &entryMessageContext{Entry: entry, param: param}
}

var _ MessageContext = (*entryMessageContext)(nil)

type entryMessageContext struct {
	Entry
	param *MessageParameters
}

// Err returns the error of this message, if set by the sender.
func (m entryMessageContext) Err() error {
	return nil
}

// ForceSend If true, Sender asks for this message to always be sent immediately.
func (m entryMessageContext) ForceSend() bool {
	return m.param.ForceSend
}

// Cooldown Gets the cooldown for this message.
func (m entryMessageContext) Cooldown() time.Duration {
	return m.param.Cooldown
}

// Maleo Gets the Maleo instance that created this MessageContext.
func (m entryMessageContext) Maleo() *Maleo {
	return m.param.Maleo
}
