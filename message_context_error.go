package maleo

import "time"

type ErrorMessageContextConstructor interface {
	BuildErrorMessageContext(err Error, param *MessageParameters) MessageContext
}

var _ ErrorMessageContextConstructor = (ErrorMessageConstructorFunc)(nil)

type ErrorMessageConstructorFunc func(err Error, param *MessageParameters) MessageContext

func (f ErrorMessageConstructorFunc) BuildErrorMessageContext(err Error, param *MessageParameters) MessageContext {
	return f(err, param)
}

func defaultErrorMessageContextConstructor(err Error, param *MessageParameters) MessageContext {
	return &errorMessageContext{Error: err, param: param}
}

var _ MessageContext = (*errorMessageContext)(nil)

type errorMessageContext struct {
	Error
	param *MessageParameters
}

// Err returns the error of this message.
func (e errorMessageContext) Err() error {
	return e.Error
}

// ForceSend If returns true, Sender asks for this message to always be sent.
func (e errorMessageContext) ForceSend() bool {
	return e.param.ForceSend
}

// Maleo Gets the maleo instance that created this MessageContext.
func (e errorMessageContext) Maleo() *Maleo {
	return e.param.Maleo
}

// Cooldown Gets the cooldown for this message.
func (e errorMessageContext) Cooldown() time.Duration {
	return e.param.Cooldown
}
