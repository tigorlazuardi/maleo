package maleo

import (
	"context"
	"time"
)

// EntryBuilder is the builder for Entry.
type EntryBuilder interface {
	// Code Sets the code for this entry.
	Code(i int) EntryBuilder

	// Message Sets the message for this entry.
	//
	// In built in implementation, If args are supplied, fmt.Sprintf will be called with s as base string.
	//
	// Very unlikely you will need to set this, because maleo already create the message field for you when you call maleo.NewEntry.
	Message(s string, args ...any) EntryBuilder

	// Context Sets additional data that will enrich how the entry will look.
	//
	// `maleo.Fields` is a type that is more integrated with built-in Messengers.
	// Using this type as Context value will often have special treatments for it.
	//
	// In built-in implementation, additional call to .Context() will make additional index, not replacing what you already set.
	//
	// Example:
	//
	// 	maleo.NewEntry(msg).Code(200).Context(maleo.F{"foo": "bar"}).Freeze()
	Context(ctx ...any) EntryBuilder

	// Key Sets the key for this entry. This is how Messenger will use to identify if an entry is the same as previous or not.
	//
	// In maleo's built-in implementation, by default, no key is set when creating new entry.
	//
	// Usually by not setting the key, The Messenger will generate their own key for this message.
	//
	// In built in implementation, If args are supplied, fmt.Sprintf will be called with key as base string.
	Key(key string, args ...any) EntryBuilder

	// Caller Sets the caller for this entry.
	//
	// In maleo's built-in implementation, by default, the caller is the location where you call `maleo.NewEntry`.
	Caller(c Caller) EntryBuilder

	// Time Sets the time for this entry. By default, it's already set when you call maleo.NewEntry.
	Time(time.Time) EntryBuilder

	// Level Sets the level for this entry.
	//
	// In maleo's built-in implementation, this defaults to what method you call to generate this entry.
	Level(lvl Level) EntryBuilder

	// Freeze this entry. Preventing further mutations.
	Freeze() Entry

	// Log this entry. Implicitly calling .Freeze() method.
	Log(ctx context.Context) Entry

	// Notify Sends this Entry to Messengers. Implicitly calling .Freeze() method.
	Notify(ctx context.Context, opts ...MessageOption) Entry
}

type Entry interface {
	CallerHint
	CodeHint
	ContextHint
	HTTPCodeHint
	KeyHint
	LevelHint
	MessageHint
	ServiceHint
	TimeHint

	/*
		Logs this entry.
	*/
	Log(ctx context.Context) Entry
	/*
		Notifies this entry to Messengers.
	*/
	Notify(ctx context.Context, opts ...MessageOption) Entry
}

type EntryConstructorContext struct {
	Caller  Caller
	Message string
	Maleo   *Maleo
}

type EntryConstructor interface {
	ConstructEntry(*EntryConstructorContext) EntryBuilder
}

var _ EntryConstructor = (EntryConstructorFunc)(nil)

type EntryConstructorFunc func(*EntryConstructorContext) EntryBuilder

func (f EntryConstructorFunc) ConstructEntry(ctx *EntryConstructorContext) EntryBuilder {
	return f(ctx)
}

func defaultEntryConstructor(ctx *EntryConstructorContext) EntryBuilder {
	return &entryBuilder{
		message: ctx.Message,
		caller:  ctx.Caller,
		context: []any{},
		level:   InfoLevel,
		maleo:   ctx.Maleo,
		time:    time.Now(),
	}
}
