package maleo

import (
	"context"
	"fmt"
	"time"
)

type entryBuilder struct {
	code    int
	message string
	caller  Caller
	context []any
	key     string
	level   Level
	time    time.Time
	maleo   *Maleo
}

func (e *entryBuilder) Code(i int) EntryBuilder {
	e.code = i
	return e
}

func (e *entryBuilder) Time(t time.Time) EntryBuilder {
	e.time = t
	return e
}

func (e *entryBuilder) Message(s string, args ...any) EntryBuilder {
	if len(args) > 0 {
		e.message = fmt.Sprintf(s, args...)
	} else {
		e.message = s
	}
	return e
}

func (e *entryBuilder) Context(ctx ...any) EntryBuilder {
	e.context = append(e.context, ctx...)
	return e
}

func (e *entryBuilder) Key(key string, args ...any) EntryBuilder {
	e.key = key
	return e
}

func (e *entryBuilder) Caller(c Caller) EntryBuilder {
	e.caller = c
	return e
}

func (e *entryBuilder) Level(lvl Level) EntryBuilder {
	e.level = lvl
	return e
}

func (e *entryBuilder) Freeze() Entry {
	return EntryNode{e}
}

func (e *entryBuilder) Log(ctx context.Context) Entry {
	return e.Freeze().Log(ctx)
}

func (e *entryBuilder) Notify(ctx context.Context, opts ...MessageOption) Entry {
	return e.Freeze().Notify(ctx, opts...)
}
