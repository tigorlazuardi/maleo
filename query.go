package maleo

import "errors"

// Query is a namespace group that holds the maleo's Query functions.
//
// Methods and functions under Query are utilities to search values in the error stack.
var Query query

type query struct{}

/*
GetHTTPCode Search for any error in the stack that implements HTTPCodeHint and return that value.

The API searches from the outermost error, and will return the first value it found.

Return 500 if there's no error that implements HTTPCodeHint in the stack.
*/
func (query) GetHTTPCode(err error) (code int) {
	if err == nil {
		return 500
	}

	if ch, ok := err.(HTTPCodeHint); ok { //nolint:errorlint
		return ch.HTTPCode()
	}

	return Query.GetHTTPCode(errors.Unwrap(err))
}

/*
GetCodeHint Search for any error in the stack that implements CodeHint and return that value.

The API searches from the outermost error, and will return the first value it found.

Return 500 if there's no error that implements CodeHint in the stack.

Used by maleo to search Code.
*/
func (query) GetCodeHint(err error) (code int) {
	if err == nil {
		return 500
	}

	if ch, ok := err.(CodeHint); ok { //nolint:errorlint
		return ch.Code()
	}

	return Query.GetCodeHint(errors.Unwrap(err))
}

/*
GetMessage Search for any error in the stack that implements MessageHint and return that value.

The API searches from the outermost error, and will return the first value it found.

Return empty string if there's no error that implements MessageHint in the stack.

Used by maleo to search Message in the error.
*/
func (query) GetMessage(err error) (message string) {
	if err == nil {
		return ""
	}

	if ch, ok := err.(MessageHint); ok { //nolint:errorlint
		return ch.Message()
	}

	return Query.GetMessage(errors.Unwrap(err))
}

/*
SearchCode Search the error stack for given code.

Given Code will be tested and the maleo.Error is returned if:

 1. Any of the error in the stack implements CodeHint interface, matches the given code, and can be cast to maleo.Error.
 2. Any of the error in the stack implements HTTPCodeHint interface, matches the given code, and can be cast to maleo.Error.

Otherwise, this function will look deeper into the stack and
eventually returns nil when nothing in the stack implements those three and have the code.

The search operation is "Breath First", meaning the maleo.Error is tested for CodeHint and HTTPCodeHint first before moving on.
*/
func (query) SearchCode(err error, code int) Error {
	if err == nil {
		return nil
	}

	// It's important to not use the other Search API for brevity.
	// This is because we are

	if ch, ok := err.(CodeHint); ok && ch.Code() == code { //nolint:errorlint
		if err, ok := err.(Error); ok { //nolint:errorlint
			return err
		}
	}

	if ch, ok := err.(HTTPCodeHint); ok && ch.HTTPCode() == code { //nolint:errorlint
		if err, ok := err.(Error); ok { //nolint:errorlint
			return err
		}
	}

	return Query.SearchCode(errors.Unwrap(err), code)
}

/*
SearchCodeHint Search the error stack for given code.

Given Code will be tested and the maleo.Error is returned if any of the error in the stack implements CodeHint interface, matches the given code, and can be cast to maleo.Error.

Otherwise, this function will look deeper into the stack and eventually returns nil when nothing in the stack implements CodeHint.
*/
func (query) SearchCodeHint(err error, code int) Error {
	if err == nil {
		return nil
	}

	if ch, ok := err.(CodeHint); ok && ch.Code() == code { //nolint:errorlint
		if err, ok := err.(Error); ok { //nolint:errorlint
			return err
		}
	}

	return Query.SearchCodeHint(errors.Unwrap(err), code)
}

/*
SearchHTTPCode Search the error stack for given code.

Given Code will be tested and the maleo.Error is returned if any of the error in the stack implements HTTPCodeHint interface, matches the given code, and can be cast to maleo.Error.

Otherwise, this function will look deeper into the stack and eventually returns nil when nothing in the stack implements HTTPCodeHint.
*/
func (query) SearchHTTPCode(err error, code int) Error {
	if err == nil {
		return nil
	}

	if ch, ok := err.(HTTPCodeHint); ok && ch.HTTPCode() == code { //nolint:errorlint
		if err, ok := err.(Error); ok { //nolint:errorlint
			return err
		}
	}

	return Query.SearchHTTPCode(errors.Unwrap(err), code)
}

// CollectErrors Collects all the maleo.Error in the error stack.
//
// It is sorted from the top most error to the bottom most error.
func (query) CollectErrors(err error) []Error {
	return collectErrors(err, make([]Error, 0, 4))
}

func collectErrors(err error, input []Error) []Error {
	if err == nil {
		return input
	}

	if err, ok := err.(Error); ok { //nolint:errorlint
		input = append(input, err)
	}

	return collectErrors(errors.Unwrap(err), input)
}

type ErrorStack struct {
	Caller Caller
	Error  error
}

// GetStack Gets the error stack by checking CallerHint.
//
// maleo recursively checks the given error if it implements CallerHint until all the error in the stack are checked.
//
// If you wish to get list of maleo.Error use CollectErrors instead.
func (query) GetStack(err error) []ErrorStack {
	in := make([]ErrorStack, 0, 8)
	return getStackList(err, in)
}

func getStackList(err error, input []ErrorStack) []ErrorStack {
	if err == nil {
		return input
	}
	if ch, ok := err.(CallerHint); ok { //nolint:errorlint
		return append(input, ErrorStack{Caller: ch.Caller(), Error: err})
	}
	return getStackList(errors.Unwrap(err), input)
}

// TopError Gets the Top most maleo.Error instance in the error stack.
// Returns nil if no maleo.Error instance found in the stack.
func (query) TopError(err error) Error {
	if err == nil {
		return nil
	}

	if e, ok := err.(Error); ok { //nolint:errorlint
		return e
	}

	return Query.TopError(errors.Unwrap(err))
}

// BottomError Gets the bottom most maleo.Error instance in the error stack.
// Returns nil if no maleo.Error instance found in the stack.
func (query) BottomError(err error) Error {
	top := Query.TopError(err)
	if top == nil {
		return nil
	}
	var result Error
	unwrapped := top.Unwrap()
	if e, ok := unwrapped.(Error); ok { //nolint:errorlint
		result = e
	} else {
		result = top
	}
	for unwrapped != nil {
		if e, ok := unwrapped.(Error); ok { //nolint:errorlint
			result = e
		}
		unwrapped = errors.Unwrap(unwrapped)
	}

	return result
}

// Cause returns the root cause.
func (query) Cause(err error) error {
	unwrapped := errors.Unwrap(err)
	for unwrapped != nil {
		err = unwrapped
		unwrapped = errors.Unwrap(err)
	}
	return err
}
