# Quick Start

**Maleo by default does nothing**. It will still do the wrapping and some nice stuffs around as one might expect, but
every call to `Log` and `Notify` from all sources will be a Noop if you don't configure it.

It is by design as it can interop easily with unit tests. So there's no need to call `maleo.SetGlobal` on unit tests.

The examples shown in this page are usable and have sensible defaults, but more likely you will need to personalize your
configuration to match what you or your team needed.

If you just want the full code example setup, you can skip directly to
[Full Code Example](#example-full-code-setup-and-usage).

## Initial Setup

`Maleo` by itself does nothing special. It collects information and data to be able to be used by other services or
libraries, but on its own, it does not have output to be consumable by humans i.e. logs and messages.

```go title="Minimum Setup"
func SetupMaleo() {
	// Not required to fill all fields,
	// but at the very least, fill the "Name"
	// for easy distinctions in the Logger and Notifications.
	service := maleo.Service{
		Name:          "my-service",
		Type:          "http-server",
		Environment:   "production",
		Version:       "v0.1.0"
	}
	mal := maleo.New(service)
	maleo.SetGlobal(mal)
}
```

`maleo.SetGlobal` function will make `*maleo.Maleo` instance you passed in as the default generator and executor to
generate entries, errors, logging, and send data to Messengers. Those will be explained later.

!!! Note "Call `maleo.SetGlobal` as early as possible in your application"

    `Maleo` does not plan to support concurrency handling on switching up global instances.
    They are unnecessary abstractions since you can always create a new instance of `*maleo.Maleo` and store them
    inside a struct or something similar.

For now, you will likely notice the snippet above **_only_** sets the metadata that would likely be very useful in
debugging your program. The code above works and compiles, but we still have not set any outputs.

```go title="Example Error Handling (after doing Setup)"
if err != nil {
	return maleo.Wrap(err).Message("failed to execute read operation").Log(ctx)
}
```

Doing operation like above will not yield any output when you call `.Log(context.Context)` method chain as expected.
However, information like `Message` and `Error Stack` are collected.

## Adding [Logger]

[Logger], by `Maleo's` definition, is whatever implements the following interface:

```go title="Logger Interface"
--8<-- "logger.go:interface"
```

`Maleo` goes out of the norm by not following the popular convention of `Info`, `Infof`, `Error`, `Errorf`, and so on.

The reason is simple, `Maleo` wants to give as much concrete information as possible to allow implementors to create
innovative, flexible, query-able, and standardized structured logs but still dynamic and developer friendly enough to
allow passing information with rich context as painless as possible.

### [Zap] as [Logger] for `Maleo`

We will use [Zap] as the [Logger] engine for the simple reason that `Maleo` has submodule `maleozap` as bridge between
[Zap] and `Maleo` for painless integration.

Run the command below to get started:

```sh
go get github.com/tigorlazuardi/maleo/maleozap
```

The command above will pull `maleozap` integration and also [Zap], since it's required by the former.

The setup is rather simple, you just have to pass on an instance of
[\*zap.Logger](https://pkg.go.dev/go.uber.org/zap#Logger) to `maleozap`.

```go title="Setup Logger"
func SetupLogger() (maleo.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Error("failed to setup zap.NewProduction(): %w", err)
	}
	mlog := maleozap.New(logger)
	return mlog, nil
}
```

!!! Info ""

    For more details and configurations check [Zap's Integration Page].

Now we just have to combine with our first code above.

```go title="Combine With Logger" linenums="1" hl_lines="18-22"
func SetupLogger() (maleo.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Error("failed to setup zap.NewProduction(): %w", err)
	}
	mlog := maleozap.New(logger)
	return mlog, nil
}

func SetupMaleo() {
	service := maleo.Service{
		Name:          "my-service",
		Type:          "http-server",
		Environment:   "production",
		Version:       "v0.1.0"
	}
	mal := maleo.New(service)
	if log, err := SetupLogger(); err != nil {
		fmt.Println(err.Error())
	} else {
		mal.SetLogger(log)
	}
	maleo.SetGlobal(mal)
}
```

Notice the highlighted lines. That's one way to register a [Logger] to Maleo.

Now let's try how the logging works.

=== "Code"

    ```go title="Obvious Error Code" hl_lines="4-7"
    func parse() error {
    	_, err = strconv.Atoi("not a number")
    	if err != nil {
    		return maleo.Wrap(err).
    			Message("failed to parse number").
    			Context(maleo.F{"input": "not a number"}).
    			Log(ctx)
    	}
    }
    ```

=== "Log Output"

    ```json title="Zap Output Using Production Config (Pretty Output for Readability)"
    {
      "level": "error",
      "ts": 1674363070.9421647,
      "caller": "zap_discord_s3/main.go:79",
      "msg": "failed to parse number",
      "service": {
    	"name": "my-service",
    	"type": "http-server",
    	"environment": "production",
    	"version": "v0.1.0"
      },
      "code": 500,
      "context": {
    	"input": "not a number"
      },
      "error": {
    	"summary": "strconv.Atoi: parsing \"not a number\": invalid syntax",
    	"details": {
    	  "Func": "Atoi",
    	  "Num": "not a number",
    	  "Err": {}
    	}
      },
      "stacktrace": "<Stack Trace>"
    }
    ```

Now we have output in `stderr`. The `service` is there for queries, and the `context` is whatever we input also there.
`context` can be used to distinguish if the users of your App and API properly follows your documentation or as proof
for your third-party that the fault are theirs. More information can be found on [Zap's Integration Page] such as why
the JSON output is like above.

## Adding Messenger

## Example Full Code Setup and Usage

```go title="_example/zap_discord/main.go" linenums="1"
--8<-- "_example/zap_discord/main.go"
```

[Logger]: ./../documentation/logger/index.md
[Zap]: https://github.com/uber-go/zap
[Zap's Integration Page]: ../documentation/logger/zap.md
