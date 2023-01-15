package maleo

import "time"

type HTTPCodeHint interface {
	// HTTPCode Gets HTTP Status Code for the type.
	HTTPCode() int
}

type CodeHint interface {
	// Code Gets the original code of the type.
	Code() int
}

type CallerHint interface {
	// Caller returns the caller of this type.
	Caller() Caller
}

type MessageHint interface {
	// Message returns the message of the type.
	Message() string
}

type KeyHint interface {
	// Key returns the key for this type.
	Key() string
}

type ContextHint interface {
	// Context returns the context of this type.
	Context() []any
}

type ServiceHint interface {
	// Service returns the service information.
	Service() Service
}

type LevelHint interface {
	// Level returns the level of this type.
	Level() Level
}

type TimeHint interface {
	// Time returns the time of this type.
	Time() time.Time
}
