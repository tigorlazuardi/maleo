package maleohttp

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/tigorlazuardi/maleo"
)

type BodyTransformer interface {
	// BodyTransform transform given input into another shape. This is called before the Encoder.
	//
	// input may be nil so make sure to take account for such situation.
	//
	// If the returned value is nil, the process will be short-circuited and the response body will be empty.
	BodyTransform(ctx context.Context, input any) any
}

type ErrorBodyTransformer interface {
	// ErrorBodyTransform transform given input into another shape. This is called before the Encoder.
	//
	// input may be nil so make sure to take account for such situation.
	//
	// If the returned value is nil, the process will be short-circuited and the response body will be empty.
	ErrorBodyTransform(ctx context.Context, err error) any
}

// BodyTransformFunc is a convenient function that implements BodyTransformer.
type BodyTransformFunc func(ctx context.Context, input any) any

func (b BodyTransformFunc) BodyTransform(ctx context.Context, input any) any {
	return b(ctx, input)
}

// NoopBodyTransform is a BodyTransformer that does nothing and only return the input as is.
type NoopBodyTransform struct{}

func (n NoopBodyTransform) BodyTransform(_ context.Context, input any) any {
	return input
}

type SimpleErrorTransformer struct{}

func (n SimpleErrorTransformer) ErrorBodyTransform(_ context.Context, err error) any {
	var msg any
	if err == nil {
		err = errors.New("[nil]")
	}
	switch err := err.(type) {
	case maleo.MessageHint:
		msg = err.Message()
	case json.Marshaler:
		msg = err
	default:
		msg = err.Error()
	}
	return map[string]any{"error": msg}
}
