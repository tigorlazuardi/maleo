# Bucket

Bucket provides a way for Messengers to store files.

You can add your own Bucket by implementing the following signature:

```go title="Bucket Interface"
type Bucket interface {
	// Upload File(s) to the bucket.
	// If File.Data() implements io.Closer, the close method should be called after upload is done.
	//
	// Whether the Upload operation is successful or not. The output must be returned.
	//
	// The number of UploadResult must be the same as the number of files uploaded.
	Upload(ctx context.Context, files []File) []UploadResult
}
```

The `Bucket.Upload` method is the entry point for you to execute your own upload logic.

Rate limit and retry are handled by the implementor (a.k.a. you).

Maleo's built-in Messengers will wait for the `Upload` method to return before sending the message.
The same cannot be said for custom Messengers. Consult the documentation of your custom Messenger for more information.