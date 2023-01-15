package maleo

import (
	"strings"
	"time"
)

type MessageParameters struct {
	ForceSend  bool
	Messengers Messengers
	Benched    Messengers
	Cooldown   time.Duration
	Maleo      *Maleo
}

func (m *MessageParameters) clone() *MessageParameters {
	mess := make(Messengers, len(m.Messengers))
	copy(mess, m.Messengers)
	bench := make(Messengers, len(m.Benched))
	copy(bench, m.Benched)
	return &MessageParameters{
		ForceSend:  m.ForceSend,
		Messengers: mess,
		Benched:    bench,
		Cooldown:   m.Cooldown,
		Maleo:      m.Maleo,
	}
}

type MessageOption interface {
	Apply(*MessageParameters)
}

type (
	MessageOptionBuilder []MessageOption
	MessageOptionFunc    func(*MessageParameters)
)

func (m MessageOptionFunc) Apply(parameters *MessageParameters) {
	m(parameters)
}

func (m MessageOptionBuilder) Apply(parameters *MessageParameters) {
	for _, opt := range m {
		opt.Apply(parameters)
	}
}

// ForceSend asks the Messengers to send the message immediately.
// Typically, the implementation of Messengers will still take account of Messenger's own rate limits,
// so the message may still be delayed or dropped. What the Messenger will do however is to ignore any message cooldowns.
//
// So if the Messenger has a rate limit of 1 second per message, ForceSend will send the Message after 1 second of last
// message sent, regardless of the cooldown.
//
// Note: Setting this option when initializing Maleo will basically ask the Messengers for *all* messages to be sent immediately.
func (m MessageOptionBuilder) ForceSend(b bool) MessageOptionBuilder {
	return append(m, MessageOptionFunc(func(p *MessageParameters) {
		p.ForceSend = b
	}))
}

// Messengers set the Messengers you want to use and ignores the rest for this message.
//
// Using this option when initializing Maleo will register the Messenger as default Messengers for Maleo.
//
// Alternatively, use Maleo.Register or Maleo.RegisterBenched to register Messengers.
func (m MessageOptionBuilder) Messengers(messengers ...Messenger) MessageOptionBuilder {
	return append(m, MessageOptionFunc(func(p *MessageParameters) {
		p.Messengers = messengers
	}))
}

// Include adds extra Messengers to the list of Messengers to be called.
func (m MessageOptionBuilder) Include(messengers ...Messenger) MessageOptionBuilder {
	return append(m, MessageOptionFunc(func(p *MessageParameters) {
		p.Messengers = append(p.Messengers, messengers...)
	}))
}

// IncludeBenched adds the Messengers that are benched to the list of Messengers to be called.
func (m MessageOptionBuilder) IncludeBenched() MessageOptionBuilder {
	return append(m, MessageOptionFunc(func(p *MessageParameters) {
		m.Include(p.Benched...)
	}))
}

// IncludeBenchedPrefix includes the benched Messengers have the given prefix in their name to be used by Maleo.
func (m MessageOptionBuilder) IncludeBenchedPrefix(prefix string) MessageOptionBuilder {
	return append(m, MessageOptionFunc(func(p *MessageParameters) {
		p.Messengers = append(p.Messengers, p.Benched.Filter(func(messenger Messenger) bool {
			return strings.HasPrefix(messenger.Name(), prefix)
		})...)
	}))
}

// IncludeBenchedSuffix includes the benched Messengers have the given suffix in their name to be used by Maleo.
func (m MessageOptionBuilder) IncludeBenchedSuffix(suffix string) MessageOptionBuilder {
	return m.IncludeBenchedFilter(func(messenger Messenger) bool {
		return strings.HasSuffix(messenger.Name(), suffix)
	})
}

// IncludeBenchedContains includes the benched Messengers that have the given string in their name to be used by Maleo.
func (m MessageOptionBuilder) IncludeBenchedContains(str string) MessageOptionBuilder {
	return m.IncludeBenchedFilter(func(messenger Messenger) bool {
		return strings.Contains(messenger.Name(), str)
	})
}

// IncludeBenchedName includes the benched Messengers that have the exact given name(s).
func (m MessageOptionBuilder) IncludeBenchedName(names ...string) MessageOptionBuilder {
	if len(names) == 0 {
		return m
	}
	return m.IncludeBenchedFilter(func(messenger Messenger) bool {
		for _, n := range names {
			if messenger.Name() == n {
				return true
			}
		}
		return false
	})
}

// IncludeBenchedFilter filters the benched Messengers to only include those that match the given filter.
func (m MessageOptionBuilder) IncludeBenchedFilter(f FilterMessengersFunc) MessageOptionBuilder {
	return m.Filter(nil, f)
}

// Exclude removes the Messengers that match the given filter. Does not affect the benched Messengers.
func (m MessageOptionBuilder) Exclude(f FilterMessengersFunc) MessageOptionBuilder {
	return append(m, MessageOptionFunc(func(p *MessageParameters) {
		p.Messengers = p.Messengers.Filter(func(messenger Messenger) bool {
			return !f(messenger)
		})
	}))
}

// ExcludePrefix removes the Messengers that match the given filter. Does not affect the benched Messengers.
func (m MessageOptionBuilder) ExcludePrefix(prefix string) MessageOptionBuilder {
	return m.Exclude(func(messenger Messenger) bool {
		return strings.HasPrefix(messenger.Name(), prefix)
	})
}

// ExcludeSuffix removes the Messengers that match the given filter. Does not affect the benched Messengers.
func (m MessageOptionBuilder) ExcludeSuffix(suffix string) MessageOptionBuilder {
	return m.Exclude(func(messenger Messenger) bool {
		return strings.HasSuffix(messenger.Name(), suffix)
	})
}

// ExcludeName removes the Messengers that have the exact given name(s). Does not affect the benched Messengers.
func (m MessageOptionBuilder) ExcludeName(name ...string) MessageOptionBuilder {
	if len(name) == 0 {
		return m
	}
	return m.Exclude(func(messenger Messenger) bool {
		for _, n := range name {
			// Exclude function reverses the result of the filter function.
			// A value of true means the Messenger is excluded.
			if messenger.Name() == n {
				return true
			}
		}
		return false
	})
}

// Filter filters the Messengers to only include those that match the given filter.
//
// It's safe to use nil for either or both arguments, if you don't wish to filter. They are merely ignored when nil.
//
// The result of two filters are then combined for Maleo to call.
func (m MessageOptionBuilder) Filter(filterMessenger, filterBenched FilterMessengersFunc) MessageOptionBuilder {
	return append(m, MessageOptionFunc(func(p *MessageParameters) {
		if filterMessenger != nil {
			p.Messengers = p.Messengers.Filter(filterMessenger)
		}
		if filterBenched != nil {
			p.Messengers = append(p.Messengers, p.Benched.Filter(filterBenched)...)
		}
	}))
}

// FilterName filters the Messengers to only include those that have the exact given name(s).
//
// This option also include Messengers that are benched into the list of Messengers to be called if they match the name.
//
// This option is ignored when you give empty names.
func (m MessageOptionBuilder) FilterName(names ...string) MessageOptionBuilder {
	if len(names) == 0 {
		return m
	}
	f := func(messenger Messenger) bool {
		for _, name := range names {
			if messenger.Name() == name {
				return true
			}
		}
		return false
	}
	return m.Filter(f, f)
}

// Cooldown overrides the base cooldown for the message for the Messengers.
//
// The implementation differs per Messenger and may instead ignore this value completely. But the builtin Messengers
// will use this value as the base cooldown before multiplying it with the cooldown multiplier for every same message
// sent.
//
// e.g. setting this cooldown to 5 seconds and the cooldown multiplier to 2 will result in the following cooldowns:
//
// 1st message: 5 seconds
//
// 2nd message: 10 seconds
//
// 3rd message: 20 seconds
//
// 4th message: 40 seconds
//
// and so on.
//
// Setting this option to 0 or negative value will result in the Messenger using the default cooldown.
//
// Setting this option to initialize Maleo will set the default cooldown for all messages.
func (m MessageOptionBuilder) Cooldown(d time.Duration) MessageOptionBuilder {
	return append(m, MessageOptionFunc(func(p *MessageParameters) {
		if d <= 0 {
			p.Cooldown = time.Minute * 15
			return
		}
		p.Cooldown = d
	}))
}
