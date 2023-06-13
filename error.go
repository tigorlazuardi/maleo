package maleo

import (
	"context"
	"fmt"
	"time"
)

type ErrorConstructorContext struct {
	Err            error
	Caller         Caller
	MessageAndArgs []any
	Maleo          *Maleo
}

type ErrorConstructor interface {
	ConstructError(*ErrorConstructorContext) ErrorBuilder
}

type ErrorConstructorFunc func(*ErrorConstructorContext) ErrorBuilder

func (f ErrorConstructorFunc) ConstructError(ctx *ErrorConstructorContext) ErrorBuilder {
	return f(ctx)
}

func defaultErrorGenerator(ctx *ErrorConstructorContext) ErrorBuilder {
	var message string
	if len(ctx.MessageAndArgs) > 0 {
		var fmtMessage string
		if msg, ok := ctx.MessageAndArgs[0].(string); ok {
			fmtMessage = msg
		} else {
			fmtMessage = fmt.Sprint(ctx.MessageAndArgs[0])
		}
		if len(ctx.MessageAndArgs) > 1 {
			message = fmt.Sprintf(fmtMessage, ctx.MessageAndArgs[1:]...)
		} else {
			message = fmtMessage
		}
	} else {
		if msg := Query.GetMessage(ctx.Err); msg != "" {
			message = msg
		} else {
			message = ctx.Err.Error()
		}
	}
	return &errorBuilder{
		code:    Query.GetCodeHint(ctx.Err),
		message: message,
		caller:  ctx.Caller,
		context: []any{},
		level:   ErrorLevel,
		origin:  ctx.Err,
		maleo:   ctx.Maleo,
		time:    time.Now(),
	}
}

/*
Error is an interface providing read only values to the error, and because it's read only, this is safe for multithreaded use.
*/
type Error interface {
	error
	CallerHint
	CodeHint
	ContextHint
	UnwrapError
	HTTPCodeHint
	KeyHint
	ErrorWriter
	LevelHint
	MessageHint
	TimeHint
	ServiceHint

	// Log this error.
	Log(ctx context.Context) Error
	// Notify notifies this error to Messengers.
	Notify(ctx context.Context, opts ...MessageOption) Error
}

type UnwrapError interface {
	// Unwrap Returns the error that is wrapped by this error. To be used by errors.Is and errors.As functions from errors library.
	Unwrap() error
}
