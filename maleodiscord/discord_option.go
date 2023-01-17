package maleodiscord

import (
	"time"

	"github.com/tigorlazuardi/maleo"
	"github.com/tigorlazuardi/maleo/bucket"
	"github.com/tigorlazuardi/maleo/locker"
)

type DiscordOption interface {
	apply(*Discord)
}

type discordOptionFunc func(*Discord)

func (d discordOptionFunc) apply(discord *Discord) {
	d(discord)
}

// WithName sets the name of this messenger.
func WithName(name string) DiscordOption {
	return discordOptionFunc(func(discord *Discord) {
		discord.name = name
	})
}

// WithLock sets the locker engine.
func WithLock(lock locker.Locker) DiscordOption {
	return discordOptionFunc(func(discord *Discord) {
		discord.lock = lock
	})
}

// WithSemaphore sets the number of concurrent workers.
func WithSemaphore(sem chan struct{}) DiscordOption {
	return discordOptionFunc(func(discord *Discord) {
		discord.sem = sem
	})
}

// WithTrace sets the tracer.
func WithTrace(trace maleo.TraceCapturer) DiscordOption {
	return discordOptionFunc(func(discord *Discord) {
		discord.trace = trace
	})
}

func WithEmbedBuilder(builder EmbedBuilder) DiscordOption {
	return discordOptionFunc(func(discord *Discord) {
		discord.builder = builder
	})
}

func WithBucket(bucket bucket.Bucket) DiscordOption {
	return discordOptionFunc(func(discord *Discord) {
		discord.bucket = bucket
	})
}

func WithGlobalKey(key string) DiscordOption {
	return discordOptionFunc(func(discord *Discord) {
		discord.globalKey = key
	})
}

func WithCooldown(cooldown time.Duration) DiscordOption {
	return discordOptionFunc(func(discord *Discord) {
		discord.cooldown = cooldown
	})
}

func WithDataEncoder(dataEncoder DataEncoder) DiscordOption {
	return discordOptionFunc(func(discord *Discord) {
		discord.dataEncoder = dataEncoder
	})
}

func WithCodeBlockBuilder(codeBlockBuilder CodeBlockBuilder) DiscordOption {
	return discordOptionFunc(func(discord *Discord) {
		discord.codeBlockBuilder = codeBlockBuilder
	})
}

func WithHook(hook Hook) DiscordOption {
	return discordOptionFunc(func(discord *Discord) {
		discord.hook = hook
	})
}

func WithClient(client Client) DiscordOption {
	return discordOptionFunc(func(discord *Discord) {
		discord.client = client
	})
}
