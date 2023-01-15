package maleo

import "context"

type Messenger interface {
	// Name Returns the name of the Messenger.
	Name() string
	// SendMessage Send Message to Messengers.
	SendMessage(ctx context.Context, msg MessageContext)

	// Wait Waits until all message in the queue or until given channel is received.
	//
	// Implementer must exit the function as soon as possible when this ctx is canceled.
	Wait(ctx context.Context) error
}

type Messengers []Messenger

type FilterMessengersFunc = func(Messenger) bool

// Filter Filters the Messengers.
func (m Messengers) Filter(f FilterMessengersFunc) Messengers {
	res := make(Messengers, 0, len(m))
	for _, messenger := range m {
		if f(messenger) {
			res = append(res, messenger)
		}
	}
	return res
}

// SendMessage Send Message to Messengers.
func (m Messengers) SendMessage(ctx context.Context, msg MessageContext) {
	for _, messenger := range m {
		messenger.SendMessage(ctx, msg)
	}
}

// Wait Waits until all message in the queue or until given channel is received.
func (m Messengers) Wait(ctx context.Context) error {
	for _, messenger := range m {
		if err := messenger.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}
