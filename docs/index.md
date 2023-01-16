# Introduction

Maleo is an _opiniated_ Golang library / framework to handle Errors, Logging, and Notification. All of those are handled
in one swoop to enhance developer experience.

Maleo's purpose is to give certain understanding on how an Error happened, but not just for the developer him/herself,
but also the team, QA, or anyone else in the team who have the same privilege and collective interest (e.g. your tech
lead or product manager who wants to understand where the fault lies).

So basically, when something goes down, the whole gang knows there's a problem and the cause of the reason is available
to them, not just the devs.

If you want to know the motivation why this exist check [here](./trivia/why-does-this-library-exist.md).

## Features

1. Designed for error chaining. Auto inference support when wrapping errors to reduce tediousness. See
   [Auto Inference](./features/auto-inference.md) for details on how it works.

    ```go title="Returning Rich Error Example"
    if err != nil {
        // Tries to detect your error when maleo.Wrap is called.
        // Filling `code` and `message` automatically when found.
        //
        // Or override them by configuring the builder that is returned by Wrap.
        return maleo.Wrap(err).
            // Code(400).
            // Context(maleo.F{"user": user}).
            // Message("failed to find user id from database").
            Freeze()
    }
    ```

2. Easy logger and notification call in one flow.

    ```go title="Easy Logger and Notification"
    if err != nil {
        return maleo.Wrap(err).
            Code(400).
            Context(maleo.F{"foo": "bar"}).
            Message("message test").
            Log(ctx).
            Notify(ctx)
    }
    ```

3. Collect only relevant Stack traces informations. While `runtime.Stack` method is available, it prints too many
   informations. Most of the time I just want to know "who calls this?", and want just that information and don't care
   about other libraries.

4. Easily extensible. You can easily add your own Logger and Messengers to use your favorite logger library and
   notification services respectively.

5. Support integration with popular libraries and platforms. Like [zap](https://github.com/uber-go/zap),
   [Amazon S3](https://aws.amazon.com/s3/), [Minio](https://min.io/), [Discord](https://discord.com/).
