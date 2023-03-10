package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"go.uber.org/zap"

	"github.com/tigorlazuardi/maleo"
	"github.com/tigorlazuardi/maleo/loader"
	"github.com/tigorlazuardi/maleo/maleodiscord"
	"github.com/tigorlazuardi/maleo/maleozap"
)

func parseInt(s string) (i int, err error) {
	i, err = strconv.Atoi(s)
	if err != nil {
		ctx := context.Background()
		return i, maleo.Wrap(err).
			Message("failed to parse '%s' into int", s).
			Log(ctx).
			Notify(ctx)
	}
	return i, err
}

func main() {
	loader.LoadEnv()
	var err error
	defer func() {
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}()
	if os.Getenv("DISCORD_WEBHOOK") == "" {
		err = fmt.Errorf("environment variable DISCORD_WEBHOOK is not set")
		return
	}

	discord := maleodiscord.NewDiscordBot(os.Getenv("DISCORD_WEBHOOK"))
	zlog, err := zap.NewProduction()
	if err != nil {
		return
	}
	zapLogger := maleozap.New(zlog)

	mal := maleo.New(maleo.Service{
		Name:        "my-service",
		Type:        "http-server",
		Version:     "v1.0.0",
		Environment: "production",
	}, maleo.Option.Init().
		Logger(zapLogger).
		Messengers(discord).
		// set caller depth to 3 so the caller function will point
		// where calling maleo.Wrap or maleo.NewEntry (exported functions) is called.
		//
		// This CallerDepth is for wrapping errors, not for logging.
		// You have to configure the logger itself to configure that.
		CallerDepth(3),
	)

	// sets `maleo.Wrap`, `maleo.NewEntry`, `maleo.Bail` functions and their variants
	// to use the new instance we just created.
	maleo.SetGlobal(mal)

	ctx := context.Background()

	// create entry and send them to discord.
	maleo.NewEntry("hello world").Log(ctx).Notify(ctx)

	// Test the error message.
	_, _ = parseInt("hello")

	// wait for discord to finish sending messages.
	err = maleo.Wait(ctx)
}
