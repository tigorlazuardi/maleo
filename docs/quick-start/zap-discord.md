---
hide:
    - toc
---

# Zap Discord Setup

```go title="Setup"
func setupMaleo() {
	// Service provides metadata to know where does this log
	// or notification comes from.
	//
	// It does not have any other impact as far as
	// Maleo and Maleo's built-in integration concerns.
	//
	// But You can use this information for your custom implementations later.
	//
	// No field is required, but it's recommended
	// to set at the very least the name and environment.
	// So when you receive a notification,
	// you can easily distinguish if it comes from prod or dev for example.
	service := maleo.Service{
		Name: "secret service",
		Environment: "development",
		Type: "kafka-consumer",
	}

	mal := maleo.NewMaleo(service, maleo.Option.Init().
		Logger(). // Set a logger to use.
		Messengers(), // Set your messengers.
	)

	maleo.SetGlobal(mal) // Set the global Maleo instance. Optional.
}
```

**_Setup is best done in the earliest possible in your code_**. `maleo.SetGlobal` does not support `Mutex` and may cause
unexpected result when called in concurrent manner.

After you finish setup Maleo, you can use it anywhere in your code.

```go title="Use"
func foo() {
	data, err := doSomethingCool()
	if err != nil {
		// Uses global instance when using direct functions.
		return nil, maleo.Wrap(err).Message("something went wrong").Log(ctx)
	}
}
```

## Testing

If you wish to test the log output, Maleo supports a very basic JSON logging outputs. It's not meant for production use.

```go title="Test"
func TestSomeFunc(t *testing.T) {
	mal, log := maleo.NewTestingMaleo()

	_ = mal.Wrap(errors.New("foo")).Log(context.Background())

	out := log.String() // or log.Bytes()
	if !strings.Contains(out, "foo") {
		t.Fatal("expected to contain foo")
	}
}
```

This method is useful when you want test implementation against Maleo. Very nice when combined with
[jsonassert](https://github.com/kinbiko/jsonassert) testing library.

Obviously, this method is not very effective if your codebase uses global instance. You have to meddle with global
instance which most likely need some further setups of your test, and thus a not recommended approach.

It's best to just test the error directly instead of the log output in this case.
