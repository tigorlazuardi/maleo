package maleo

// Wrap the error. The returned ErrorBuilder may be appended with values.
// Call .Freeze() method to turn this into proper error.
// Or call .Log() or .Notify() to implicitly freeze the error and do actual stuffs.
//
// Example:
//
//	if err != nil {
//	  return maleo.Wrap(err).Message("something went wrong").Freeze()
//	}
//
// Example with Log:
//
//	if err != nil {
//	  return maleo.Wrap(err).Message("something went wrong").Log(ctx)
//	}
//
// Example with Notify:
//
//	if err != nil {
//	  return maleo.Wrap(err).Message("something went wrong").Notify(ctx)
//	}
//
// Example with Notify and Log:
//
//	if err != nil {
//	  return maleo.Wrap(err).Message("something went wrong").Log(ctx).Notify(ctx)
//	}
func Wrap(err error) ErrorBuilder {
	return Global().Wrap(err)
}

// Bail creates a new ErrorBuilder from simple messages.
//
// If args are not empty, msg will be fed into fmt.Errorf along with the args.
// Otherwise, msg will be fed into `errors.New()`.
func Bail(msg string, args ...any) ErrorBuilder {
	return Global().Bail(msg, args...)
}

// WrapFreeze is a Shorthand for `maleo.Wrap(err).Message(message, args...).Freeze()`
//
// Useful when just wanting to add extra simple messages to the error chain.
func WrapFreeze(err error, message string, args ...any) Error {
	return Global().WrapFreeze(err, message, args...)
}

// BailFreeze creates new immutable Error from simple messages.
//
// It's a shorthand for `maleo.Bail(msg, args...).Freeze()`.
func BailFreeze(msg string, args ...any) Error {
	return Global().BailFreeze(msg, args...)
}
