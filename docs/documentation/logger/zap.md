# Zap

`maleozap` checks if an error, when fed into `json.Marshal` returns an empty bracket pair `{}`. This is common in
`fmt.Errorf` and `errors.New` errors. So `maleozap` tries to put `error.Error()` output in to `summary`field.

However, marshaling errors sometimes have strong benefits that sometimes a mere `summary` do not cover. Built-in `http`
and `pq (postgres)` driver often do this. So `maleozap` will take both to bring as much information as it can, but it
will omit `details` field if they are just empty information since they are merely a waste of space.
