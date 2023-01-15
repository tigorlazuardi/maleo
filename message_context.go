package maleo

import "time"

// MessageContext is the context of a message.
//
// It holds the message and data that can be sent to the Messenger's target.
type MessageContext interface {
	HTTPCodeHint
	CodeHint
	MessageHint
	CallerHint
	KeyHint
	LevelHint
	ServiceHint
	ContextHint
	TimeHint
	// Err returns the error item. May be nil if message contains no error.
	Err() error
	// ForceSend If true, Sender asks for this message to always be send at the earliest possible.
	ForceSend() bool
	// Cooldown returns non-zero value if Sender asks for this message to be sent after this duration.
	Cooldown() time.Duration
	// Maleo Gets the Maleo instance that created this MessageContext.
	Maleo() *Maleo
}
