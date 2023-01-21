package maleo

func NewEntry(msg string, args ...any) EntryBuilder {
	return Global().NewEntry(msg, args...)
}
