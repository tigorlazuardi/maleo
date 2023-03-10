package main

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/tigorlazuardi/maleo"
	"github.com/tigorlazuardi/maleo/bucket/maleos3-v2"
	"github.com/tigorlazuardi/maleo/loader"
	"github.com/tigorlazuardi/maleo/maleodiscord"
	"github.com/tigorlazuardi/maleo/maleozap"
)

func checkEnvs(envs ...string) error {
	for _, env := range envs {
		if os.Getenv(env) == "" {
			return fmt.Errorf("environment variable %s is not set", env)
		}
	}
	return nil
}

func main() {
	loader.LoadEnv() // load .env files.
	var err error
	defer func() {
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}()
	err = checkEnvs(
		"DISCORD_WEBHOOK",
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_KEY_ID",
		"AWS_ENDPOINT",
	)
	if err != nil {
		return
	}

	s3bucket, err := maleos3.NewS3Bucket(os.Getenv("AWS_ENDPOINT"))
	if err != nil {
		return
	}

	discord := maleodiscord.NewDiscordBot(os.Getenv("DISCORD_WEBHOOK"), maleodiscord.WithBucket(s3bucket))
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

	// wait for discord to finish sending messages.
	err = maleo.Wait(ctx)
}
