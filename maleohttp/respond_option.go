package maleohttp

import (
	"github.com/tigorlazuardi/maleo"
)

type RespondOption interface {
	Apply(*RespondContext)
}

type (
	RespondOptionBuilder []RespondOption
	RespondOptionFunc    func(*RespondContext)
)

func (r RespondOptionBuilder) Apply(o *RespondContext) {
	for _, v := range r {
		v.Apply(o)
	}
}

func (r RespondOptionFunc) Apply(o *RespondContext) {
	r(o)
}

type RespondContext struct {
	Encoder              Encoder
	BodyTransformer      BodyTransformer
	Compressor           Compressor
	StreamCompressor     StreamCompressor
	StatusCode           int
	ErrorBodyTransformer ErrorBodyTransformer
	CallerDepth          int
	Caller               maleo.Caller
}

// Encoder overrides the Encoder to be used for encoding the response body.
func (r RespondOptionBuilder) Encoder(encoder Encoder) RespondOptionBuilder {
	return append(r, RespondOptionFunc(func(o *RespondContext) {
		o.Encoder = encoder
	}))
}

// Transformer overrides the transformer to be used for transforming the response body.
func (r RespondOptionBuilder) Transformer(transformer BodyTransformer) RespondOptionBuilder {
	return append(r, RespondOptionFunc(func(o *RespondContext) {
		o.BodyTransformer = transformer
	}))
}

func (r RespondOptionBuilder) ErrorTransformer(transformer ErrorBodyTransformer) RespondOptionBuilder {
	return append(r, RespondOptionFunc(func(o *RespondContext) {
		o.ErrorBodyTransformer = transformer
	}))
}

// Compressor overrides the Compressor to be used for compressing the response body.
func (r RespondOptionBuilder) Compressor(compressor Compressor) RespondOptionBuilder {
	return append(r, RespondOptionFunc(func(o *RespondContext) {
		o.Compressor = compressor
	}))
}

// StreamCompressor overrides the StreamCompressor to be used for compressing the response body.
func (r RespondOptionBuilder) StreamCompressor(compressor StreamCompressor) RespondOptionBuilder {
	return append(r, RespondOptionFunc(func(o *RespondContext) {
		o.StreamCompressor = compressor
	}))
}

// StatusCode overrides the status code to be used for the response.
func (r RespondOptionBuilder) StatusCode(code int) RespondOptionBuilder {
	return append(r, RespondOptionFunc(func(o *RespondContext) {
		o.StatusCode = code
	}))
}

// CallerSkip overrides the caller skip to be used for the response to get the caller information.
func (r RespondOptionBuilder) CallerSkip(i int) RespondOptionBuilder {
	return append(r, RespondOptionFunc(func(o *RespondContext) {
		o.CallerDepth = i
	}))
}

// AddCallerSkip adds the caller skip value to be used for the response to get the caller information.
func (r RespondOptionBuilder) AddCallerSkip(i int) RespondOptionBuilder {
	return append(r, RespondOptionFunc(func(o *RespondContext) {
		o.CallerDepth += 1
	}))
}
