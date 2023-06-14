package maleo

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Maleo struct {
	service       Service
	defaultParams *MessageParameters
	logger        Logger
	engine        Engine
	callerDepth   int
	name          string
	isGlobal      bool
}

// New creates a new Maleo instance.
func New(service Service, opts ...InitOption) *Maleo {
	m := &Maleo{
		service: service,
		defaultParams: &MessageParameters{
			Messengers: []Messenger{},
			Benched:    []Messenger{},
			Cooldown:   time.Minute * 15,
			ForceSend:  false,
		},
		engine:      NewEngine(),
		logger:      NoopLogger{},
		callerDepth: 2,
	}
	m.defaultParams.Maleo = m
	for _, opt := range opts {
		opt.apply(m)
	}
	return m
}

// Name implements the Messenger interface. If name is not set on initialization, it will return Service.String().
func (m *Maleo) Name() string {
	if m.name == "" {
		return m.service.String()
	}
	return m.name
}

// SendMessage implements the Messenger interface. Maleo instance itself can be a Messenger for other Maleo instances.
//
// Note that SendMessage for Maleo only support Messengers that are set in initialization or call Maleo.Register.
func (m *Maleo) SendMessage(ctx context.Context, msg MessageContext) {
	for _, v := range m.defaultParams.Messengers {
		v.SendMessage(ctx, msg)
	}
}

type multierror []error

func (m multierror) Error() string {
	s := strings.Builder{}
	for i, err := range m {
		if i > 0 {
			s.WriteString("; ")
		}
		s.WriteString(strconv.Itoa(i + 1))
		s.WriteString(". ")
		s.WriteString(err.Error())
	}
	return s.String()
}

// Wait implements the Messenger interface. Maleo instance itself can be a Messenger.
func (m *Maleo) Wait(ctx context.Context) error {
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	wg.Add(len(m.defaultParams.Messengers))
	errs := make(multierror, 0, len(m.defaultParams.Messengers))
	for _, v := range m.defaultParams.Messengers {
		go func(messenger Messenger) {
			defer wg.Done()
			err := messenger.Wait(ctx)
			if err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("failed on waiting messages to finish from '%s': %w", messenger.Name(), err))
				mu.Unlock()
			}
		}(v)
	}
	wg.Wait()
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// Register registers a new messenger to the default parameters.
func (m *Maleo) Register(messengers ...Messenger) {
	m.defaultParams.Messengers = append(m.defaultParams.Messengers, messengers...)
}

// RegisterBenched registers a new messenger to the default benched parameters.
//
// Benched messengers will not be used on normal notify calls. They have to be called explicitly using
// .Include() or .Filter() options.
func (m *Maleo) RegisterBenched(messengers ...Messenger) {
	m.defaultParams.Benched = append(m.defaultParams.Benched, messengers...)
}

// Wrap works like exported maleo.Wrap, but at the scope of this Maleo instance instead.
func (m *Maleo) Wrap(err error, msgAndArgs ...any) ErrorBuilder {
	if err == nil {
		err = ErrNil
	}
	caller := GetCaller(m.callerDepth)
	return m.engine.ConstructError(&ErrorConstructorContext{
		Err:            err,
		Caller:         caller,
		Maleo:          m,
		MessageAndArgs: msgAndArgs,
	})
}

// WrapFreeze works like exported maleo.WrapFreeze, but at the scope of this Maleo instance instead.
func (m *Maleo) WrapFreeze(err error, msgAndArgs ...any) Error {
	if err == nil {
		err = ErrNil
	}
	caller := GetCaller(m.callerDepth)
	return m.engine.ConstructError(&ErrorConstructorContext{
		Err:            err,
		Caller:         caller,
		Maleo:          m,
		MessageAndArgs: msgAndArgs,
	}).Freeze()
}

// NewEntry Creates a new EntryBuilder. The returned EntryBuilder may be appended with values.
func (m *Maleo) NewEntry(msg string, args ...any) EntryBuilder {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	caller := GetCaller(m.callerDepth)
	return m.engine.ConstructEntry(&EntryConstructorContext{
		Caller:  caller,
		Maleo:   m,
		Message: msg,
	})
}

// Bail creates a new ErrorBuilder from simple messages.
//
// If args is not empty, msg will be fed into fmt.Errorf along with the args.
// Otherwise, msg will be fed into `errors.New()`.
func (m *Maleo) Bail(msg string, args ...any) ErrorBuilder {
	var err error
	if len(args) > 0 {
		err = fmt.Errorf(msg, args...)
	} else {
		err = errors.New(msg)
	}
	caller := GetCaller(m.callerDepth)
	return m.engine.ConstructError(&ErrorConstructorContext{
		Err:    err,
		Caller: caller,
		Maleo:  m,
	})
}

// BailFreeze creates new immutable Error from simple messages.
//
// If args is not empty, msg will be fed into fmt.Errorf along with the args.
// Otherwise, msg will be fed into `errors.New()`.
func (m *Maleo) BailFreeze(msg string, args ...any) Error {
	var err error
	if len(args) > 0 {
		err = fmt.Errorf(msg, args...)
	} else {
		err = errors.New(msg)
	}
	caller := GetCaller(m.callerDepth)
	return m.engine.ConstructError(&ErrorConstructorContext{
		Err:    err,
		Caller: caller,
		Maleo:  m,
	}).Freeze()
}

// Notify Sends the Entry to Messengers.
func (m *Maleo) Notify(ctx context.Context, entry Entry, parameters ...MessageOption) {
	opts := m.defaultParams.clone()
	for _, v := range parameters {
		v.Apply(opts)
	}
	msg := m.engine.BuildEntryMessageContext(entry, opts)
	m.sendNotif(ctx, msg, opts)
}

// NotifyError sends the Error to Messengers.
func (m *Maleo) NotifyError(ctx context.Context, err Error, parameters ...MessageOption) {
	opts := m.defaultParams.clone()
	for _, v := range parameters {
		v.Apply(opts)
	}
	msg := m.engine.BuildErrorMessageContext(err, opts)
	m.sendNotif(ctx, msg, opts)
}

func (m *Maleo) sendNotif(ctx context.Context, msg MessageContext, opts *MessageParameters) {
	ctx = DetachedContext(ctx)
	for _, v := range opts.Messengers {
		v.SendMessage(ctx, msg)
	}
}

func (m *Maleo) SetLogger(logger Logger) {
	m.logger = logger
}

func (m *Maleo) SetEngine(engine Engine) {
	m.engine = engine
}

// Clone creates a new Maleo instance with the same parameters as the current one.
//
// The new Maleo instance will have the same name, logger, engine, and caller depth, but marked as non-global.
func (m *Maleo) Clone() *Maleo {
	return &Maleo{
		service:       m.service,
		defaultParams: m.defaultParams.clone(),
		logger:        m.logger,
		engine:        m.engine,
		callerDepth:   m.callerDepth,
		name:          m.name,
		isGlobal:      false,
	}
}

// Log implements the Logger interface. Maleo instance itself can be a Logger for other Maleo instance.
//
// If ctx contains a Maleo instance, and current instance is a Global one, it will be used instead.
func (m *Maleo) Log(ctx context.Context, entry Entry) {
	if m.isGlobal {
		if mctx := MaleoFromContext(ctx); mctx != nil {
			mctx.Log(ctx, entry)
			return
		}
	}
	m.logger.Log(ctx, entry)
}

// LogError implements the Logger interface. Maleo instance itself can be a Logger for other Maleo instance.
//
// If ctx contains a Maleo instance, and current instance is a Global one, it will be used instead.
func (m *Maleo) LogError(ctx context.Context, err Error) {
	if m.isGlobal {
		if mctx := MaleoFromContext(ctx); mctx != nil {
			mctx.LogError(ctx, err)
			return
		}
	}
	m.logger.LogError(ctx, err)
}

func (m *Maleo) Service() Service {
	return m.service
}
