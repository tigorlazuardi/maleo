package maleo

import (
	"reflect"
	"runtime"
	"testing"
)

func TestNewEngine(t *testing.T) {
	var (
		mockErrorConstructor ErrorConstructor = ErrorConstructorFunc(func(constructorContext *ErrorConstructorContext) ErrorBuilder {
			return nil
		})
		mockEntryConstructor EntryConstructor = EntryConstructorFunc(func(constructorContext *EntryConstructorContext) EntryBuilder {
			return nil
		})
		mockEntryMessageContextConstructor EntryMessageContextConstructor = MessageContextConstructorFunc(func(entry Entry, param *MessageParameters) MessageContext {
			return nil
		})
		mockErrorMessageContextConstructor ErrorMessageContextConstructor = ErrorMessageConstructorFunc(func(err Error, param *MessageParameters) MessageContext {
			return nil
		})
	)

	opts := Option.Engine().
		ErrorMessageContextConstructor(mockErrorMessageContextConstructor).
		EntryConstructor(mockEntryConstructor).
		ErrorConstructor(mockErrorConstructor).
		EntryMessageContextConstructor(mockEntryMessageContextConstructor)

	e := NewEngine(opts)
	if e == nil {
		t.Fatalf("Expected engine to be not nil")
	}
	eng := e.(*engine)
	compareFunction(t, mockErrorConstructor, eng.ErrorConstructor)
	compareFunction(t, mockEntryConstructor, eng.EntryConstructor)
	compareFunction(t, mockEntryMessageContextConstructor, eng.EntryMessageContextConstructor)
	compareFunction(t, mockErrorMessageContextConstructor, eng.ErrorMessageContextConstructor)
}

func compareFunction(t *testing.T, expected, actual interface{}) {
	expect := runtime.FuncForPC(reflect.ValueOf(expected).Pointer()).Name()
	actualFunc := runtime.FuncForPC(reflect.ValueOf(actual).Pointer()).Name()
	if expect != actualFunc {
		t.Errorf("Expected %v, got %v", expect, actualFunc)
	}
}
