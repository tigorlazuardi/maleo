package maleo

// Wrap the error. The returned ErrorBuilder may be appended with values.
// Call .Freeze() method to turn this into proper error.
// Or call .Log() or .Notify() to implicitly freeze the error and do actual stuffs.
//
// NOTE: When you call .Log(ctx) or .Notify(ctx) and context is attached with a Maleo instance,
// the error will be logged or notified using that Maleo instance, instead of the global Maleo instance.
//
// Example:
//
//	if err != nil {
//	  return maleo.Wrap(err, "something went wrong: %s", reason).Freeze()
//	}
//
// Example with Log:
//
//	if err != nil {
//	  return maleo.Wrap(err, "something went wrong: %s", reason).Log(ctx)
//	}
//
// Example with Notify:
//
//	if err != nil {
//	  return maleo.Wrap(err, "something went wrong").Notify(ctx)
//	}
//
// Example with Notify and Log:
//
//	if err != nil {
//	  return maleo.Wrap(err, "something went wrong").Log(ctx).Notify(ctx)
//	}
//
// Example with inline Maleo instance, instead of using global:
//
//	myCustomMaleo := maleo.New()
//	ctx = maleo.ContextWithMaleo(ctx, myCustomMaleo)
//	if err != nil {
//	  return maleo.Wrap(err, "something went wrong").Log(ctx).Notify(ctx)
//	}
func Wrap(err error, msgAndArgs ...any) ErrorBuilder {
	return Global().Wrap(err, msgAndArgs...)
}

// Bail creates a new ErrorBuilder from simple messages.
//
// If args are not empty, msg will be fed into fmt.Errorf along with the args.
// Otherwise, msg will be fed into `errors.New()`.
//
// NOTE: When you call .Log(ctx) or .Notify(ctx) and context is attached with a Maleo instance,
// the error will be logged or notified using that Maleo instance, instead of the global Maleo instance.
func Bail(msg string, args ...any) ErrorBuilder {
	return Global().Bail(msg, args...)
}

// WrapFreeze is a Shorthand for `maleo.Wrap(err).Message(message, args...).Freeze()`
//
// Useful when just wanting to add extra simple messages to the error chain.
//
// NOTE: When you call .Log(ctx) or .Notify(ctx) and context is attached with a Maleo instance,
// the error will be logged or notified using that Maleo instance, instead of the global Maleo instance.
func WrapFreeze(err error, msgAndArgs ...any) Error {
	return Global().WrapFreeze(err, msgAndArgs...)
}

// BailFreeze creates new immutable Error from simple messages.
//
// It's a shorthand for `maleo.Bail(msg, args...).Freeze()`.
//
// NOTE: When you call .Log(ctx) or .Notify(ctx) and context is attached with a Maleo instance,
// the error will be logged or notified using that Maleo instance, instead of the global Maleo instance.
func BailFreeze(msg string, args ...any) Error {
	return Global().BailFreeze(msg, args...)
}
