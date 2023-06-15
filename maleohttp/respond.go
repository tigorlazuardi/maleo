package maleohttp

import (
	"context"
	"io"
	"net/http"

	"github.com/tigorlazuardi/maleo"
)

// Responder handles the response and writing to http.ResponseWriter.
type Responder struct {
	encoder          Encoder
	transformer      BodyTransformer
	errorTransformer ErrorBodyTransformer
	maleo            *maleo.Maleo
	compressor       Compressor
	streamCompressor StreamCompressor
	callerDepth      int
	hooks            RespondHookList
}

// NewResponder creates a new Responder instance.
//
// It has the following default values:
//
// - Encoder: JSONEncoder (encodes to JSON)
//
// - BodyTransformer: NoopBodyTransform (does nothing to whatever value you pass in)
//
// - ErrorBodyTransformer: SimpleErrorTransformer (encodes error to {"error": "message/err.Error()"}) with JSONEncoder.
// Different Encoder may have different output.
//
// - Maleo: points to the global maleo instance
//
// - Compressor: NoCompression.
func NewResponder() *Responder {
	return &Responder{
		encoder:          NewJSONEncoder(),
		transformer:      NoopBodyTransform{},
		errorTransformer: SimpleErrorTransformer{},
		maleo:            maleo.Global(),
		compressor:       NoCompression{},
		streamCompressor: NoCompression{},
		callerDepth:      2, // default to 3, which is wherever the user calls Responder.Respond() or it's derivatives.
	}
}

// SetErrorTransformer sets the ErrorBodyTransformer to be used by the Responder.
func (r *Responder) SetErrorTransformer(errorTransformer ErrorBodyTransformer) {
	r.errorTransformer = errorTransformer
}

// SetEncoder sets the Encoder to be used by the Responder.
func (r *Responder) SetEncoder(encoder Encoder) {
	r.encoder = encoder
}

// SetBodyTransformer sets the BodyTransformer to be used by the Responder.
func (r *Responder) SetBodyTransformer(transform BodyTransformer) {
	r.transformer = transform
}

// SetMaleo sets the maleo instance to be used by the Responder.
func (r *Responder) SetMaleo(t *maleo.Maleo) {
	r.maleo = t
}

// SetCompressor sets the compression to be used by the Responder.
func (r *Responder) SetCompressor(compressor Compressor) {
	r.compressor = compressor
}

// SetCallerDepth sets the caller depth to be used to get caller function by the Responder.
func (r *Responder) SetCallerDepth(depth int) {
	r.callerDepth = depth
}

// SetStreamCompressor sets the stream compression to be used by the Responder.
func (r *Responder) SetStreamCompressor(streamCompressor StreamCompressor) {
	r.streamCompressor = streamCompressor
}

func (r Responder) buildOption(statusCode int, rw http.ResponseWriter, request *http.Request, body any, opts ...RespondOption) *RespondContext {
	opt := &RespondContext{
		Encoder:              r.encoder,
		BodyTransformer:      r.transformer,
		Compressor:           r.compressor,
		StatusCode:           statusCode,
		ErrorBodyTransformer: r.errorTransformer,
		CallerDepth:          r.callerDepth,
		StreamCompressor:     r.streamCompressor,
	}
	for _, o := range opts {
		o.Apply(opt)
	}
	opt.Caller = maleo.GetCaller(opt.CallerDepth + 1)
	for _, hook := range r.hooks {
		opt = hook.BeforeRespond(opt, rw, request, body)
	}
	return opt
}

func (r Responder) buildOptionError(statusCode int, rw http.ResponseWriter, request *http.Request, errPayload error, opts ...RespondOption) *RespondContext {
	opt := &RespondContext{
		Encoder:              r.encoder,
		BodyTransformer:      r.transformer,
		Compressor:           r.compressor,
		StatusCode:           statusCode,
		ErrorBodyTransformer: r.errorTransformer,
		CallerDepth:          r.callerDepth,
		StreamCompressor:     r.streamCompressor,
	}
	for _, o := range opts {
		o.Apply(opt)
	}
	opt.Caller = maleo.GetCaller(opt.CallerDepth + 1)
	for _, hook := range r.hooks {
		opt = hook.BeforeRespondError(opt, rw, request, errPayload)
	}
	return opt
}

func (r Responder) buildOptionStream(statusCode int, rw http.ResponseWriter, request *http.Request, reader io.Reader, opts ...RespondOption) *RespondContext {
	opt := &RespondContext{
		Encoder:              r.encoder,
		BodyTransformer:      r.transformer,
		Compressor:           r.compressor,
		StatusCode:           statusCode,
		ErrorBodyTransformer: r.errorTransformer,
		CallerDepth:          r.callerDepth,
		StreamCompressor:     r.streamCompressor,
	}
	for _, o := range opts {
		o.Apply(opt)
	}
	opt.Caller = maleo.GetCaller(opt.CallerDepth + 1)
	for _, hook := range r.hooks {
		opt = hook.BeforeRespondStream(opt, rw, request, reader)
	}
	return opt
}

var requestBodyKey = struct{ key int }{777}

func clonedBodyFromContext(ctx context.Context) ClonedBody {
	body, ok := ctx.Value(requestBodyKey).(ClonedBody)
	if !ok {
		return nil
	}
	return body
}

func contextWithClonedBody(ctx context.Context, body ClonedBody) context.Context {
	return context.WithValue(ctx, requestBodyKey, body)
}
