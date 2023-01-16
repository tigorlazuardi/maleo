# Bucket

Bucket provides a way for Messengers to store files.

You can add your own Bucket by implementing the following signature:

```go
--8<-- "bucket/bucket.go:bucket"
```

The `Bucket.Upload` method is the entry point for you to execute your own upload logic.

Rate limit and retry are handled by the implementor (a.k.a. you).

Maleo's built-in Messengers will wait for the `Upload` method to return before sending the message.
The same cannot be said for custom Messengers. Consult the documentation of your custom Messenger for more information.