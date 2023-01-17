package maleodiscord

import (
	"context"
	"math/rand"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/bwmarrin/snowflake"

	"github.com/tigorlazuardi/maleo"
	"github.com/tigorlazuardi/maleo/bucket"
	"github.com/tigorlazuardi/maleo/locker"
	"github.com/tigorlazuardi/maleo/queue"
)

func init() {
	snowflake.Epoch = 1420070400000 // discord epoch
}

type Discord struct {
	name             string
	webhook          string
	lock             locker.Locker
	queue            *queue.Queue[*Job]
	sem              chan struct{}
	working          int32
	trace            maleo.TraceCapturer
	builder          EmbedBuilder
	bucket           bucket.Bucket
	globalKey        string
	cooldown         time.Duration
	snowflake        *snowflake.Node
	client           Client
	hook             Hook
	dataEncoder      DataEncoder
	codeBlockBuilder CodeBlockBuilder
}

// Name implements tower.Messenger interface.
func (d *Discord) Name() string {
	if d.name == "" {
		return "discord"
	}
	return d.name
}

// SendMessage implements tower.Messenger interface.
func (d *Discord) SendMessage(ctx context.Context, msg maleo.MessageContext) {
	d.queue.Enqueue(NewJob(ctx, msg))
	d.work()
}

// Wait implements tower.Messenger interface.
func (d *Discord) Wait(ctx context.Context) error {
	err := make(chan error)
	go func() {
		for {
			select {
			case <-ctx.Done():
				err <- ctx.Err()
				break
			default:
				if d.queue.Len() == 0 {
					err <- nil
					break
				}
				time.Sleep(time.Millisecond * 50)
			}
		}
	}()

	return <-err
}

func (d *Discord) work() {
	if atomic.CompareAndSwapInt32(&d.working, 0, 1) {
		go func() {
			for d.queue.HasNext() {
				d.sem <- struct{}{}
				kv := d.queue.Dequeue()
				go func() {
					ctx := maleo.DetachedContext(kv.Context)
					d.send(ctx, kv.Message)
					<-d.sem
				}()
			}
			atomic.StoreInt32(&d.working, 0)
		}()
	}
}

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

type Job struct {
	Context context.Context
	Message maleo.MessageContext
}

func NewJob(ctx context.Context, message maleo.MessageContext) *Job {
	return &Job{Context: ctx, Message: message}
}

// NewDiscordBot creates a new discord bot.
func NewDiscordBot(webhook string, opts ...DiscordOption) *Discord {
	d := &Discord{
		name:             "discord",
		webhook:          webhook,
		lock:             locker.NewLocalLock(),
		queue:            queue.New[*Job](500),
		sem:              make(chan struct{}, (runtime.NumCPU()/3)+2),
		trace:            maleo.NoopTraceCapturer{},
		globalKey:        "global",
		cooldown:         time.Minute * 15,
		snowflake:        generateSnowflakeNode(),
		client:           http.DefaultClient,
		hook:             NoopHook{},
		dataEncoder:      JSONDataEncoder{},
		codeBlockBuilder: JSONCodeBlockBuilder{},
	}
	d.builder = EmbedBuilderFunc(d.defaultEmbedBuilder)
	for _, opt := range opts {
		opt.apply(d)
	}
	return d
}

func generateSnowflakeNode() *snowflake.Node {
	source := rand.NewSource(time.Now().UnixNano())
	id := source.Int63()
	high := source.Int63()
	for high > 1023 {
		high >>= 1
	}
	for id > high {
		id >>= 1
	}
	node, _ := snowflake.NewNode(id)
	return node
}
