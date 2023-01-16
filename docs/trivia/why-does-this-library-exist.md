# Why does this library exist?

Errors are a part of life in a development life cycle. Naturally, everyone have their own way to handle errors. But,
what I found lacking in Go is support for rich error informations.

Generally, to easily debug an information, you will need to log the error (the output), and the input a.k.a. what data
is entered to produce such error.

```go title="common error handling pattern"
if err != nil {
    return fmt.Errorf("failed to find user '%s' from database: %w", userId, err)
}
```

Something like above is usually enough... when your function is rather simple in logic and the caller's logic is rather
straight forward. It became another story when you have a dozen error handling like this in the codebase.

Of course, anyone can make themselves an `error` instance that can carry such information or use a library.

But I always found something lacking. Rich Error is one thing, but sometimes I want that error thrown into my face. Of
course, to do that, I either have to pull another library or make my own. Not to mention, I also do not want to be
spammed by my own program when that Error appears (e.g. database down). I only need to know one error! Not many.

Also, there's an issue of collecive knowledge on a team.

It's nice if you work in one man project or old team on existing project. You can do something specialized with your
project when handling errors and everyone else in the team know the best practics to handle errors.

It suddenly became a pain when you have to onboard someone new to the project since you have to mentor how to handle
errors, how to write that rich error themselves, severity, how to send notification of a project, etc.

Wouldn't it be nice if there's a library to handle all of those tasks? Enter `Maleo` to handle majority of the haul of
said problem.
