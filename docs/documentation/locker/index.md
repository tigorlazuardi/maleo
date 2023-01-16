# Locker

Locker is an interface for Maleo's built-in messenger to synchronize state between services.

It is used to enable rate limit on the API (so your Messengers don't get banned) even when your services are
distributed.

To create your own `Locker` implementation, you have to implement the following interface:

```go
--8<-- "locker/locker.go:locker"
```

## Local Lock

Local lock is an in memory lock. It's a special lock that only applies to current runtime.

There are no synchronization between services, and thus states like backoff are not synchronized between services.

!!! info ""

    The lock is useful when you want to create simple prototype services, but in deployed environment, please use more
    sophisticated services to handle synchronization and persistence.

You can create a simple local lock by calling:

```go
lock := locker.NewLocalLock()
```
