---
hide:
    - toc
---

# Testing

## Testing

If you wish to test the log output, Maleo supports a very basic JSON logging outputs. Due to the simplicity and
naitivity of the API, it's considered not fit for production use.

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
