package maleohttp

import (
	"io"
	"net/http"

	"github.com/tigorlazuardi/maleo"
)

type baseHook struct {
	Context        *RespondContext
	Request        *http.Request
	RequestBody    ClonedBody
	ResponseStatus int
	ResponseHeader http.Header
	Maleo          *maleo.Maleo
	Error          error
}

type RespondBody struct {
	PreEncoded     any
	PostEncoded    []byte
	PostCompressed []byte
}

type RespondHookContext struct {
	*baseHook
	ResponseBody RespondBody
}

type RespondErrorBody struct {
	PreEncoded     error
	PostEncoded    []byte
	PostCompressed []byte
}

type RespondErrorHookContext struct {
	*baseHook
	ResponseBody RespondErrorBody
}

type RespondStreamBody struct {
	Value       ClonedBody
	ContentType string
}

type RespondStreamHookContext struct {
	*baseHook
	ResponseBody RespondStreamBody
}

type RespondHookList []RespondHook

func (r RespondHookList) CountMaximumRequestBodyRead(request *http.Request) int {
	var count int
	for _, hook := range r {
		accept := hook.AcceptRequestBodySize(request)
		// A hook requests all body to be read. We will read all body and stop looking at other hooks.
		if accept < 0 {
			count = accept
			break
		}
		if accept > count {
			count = accept
		}
	}
	return count
}

func (r RespondHookList) CountMaximumRespondBodyRead(contentType string, request *http.Request) int {
	var count int
	for _, hook := range r {
		accept := hook.AcceptResponseBodyStreamSize(contentType, request)
		if accept < 0 {
			count = accept
			break
		}
		if accept > count {
			count = accept
		}
	}
	return count
}

type RespondHook interface {
	AcceptRequestBodySize(r *http.Request) int
	AcceptResponseBodyStreamSize(respondContentType string, request *http.Request) int

	// BeforeRespond is called after populating the handlers and values but before doing any operation to the response writer.
	//
	// Implementer must not write any status code or body to the response writer.
	// It will be set by the library.
	//
	// Implementer is allowed to set any extra header or cookie to the response writer.
	//
	// Implementer can modify everything else. Return the ctx to continue the operation.
	BeforeRespond(ctx *RespondContext, rw http.ResponseWriter, request *http.Request, body any) *RespondContext
	// BeforeRespondError is called after populating the handlers and values but before doing any operation to the response writer.
	//
	// Implementer must not write any status code or body to the response writer.
	// It will be set by the library.
	//
	// Implementer is allowed to set any extra header or cookie to the response writer.
	//
	// Implementer can modify everything else. Return the ctx to continue the operation.
	BeforeRespondError(ctx *RespondContext, rw http.ResponseWriter, request *http.Request, err error) *RespondContext
	// BeforeRespondStream is called after populating the handlers and values but before doing any operation to the response writer.
	//
	// Implementer must not write any status code or body to the response writer.
	// It will be set by the library.
	//
	// Implementer is allowed to set any extra header or cookie to the response writer.
	//
	// Implementer can modify everything else. Return the ctx to continue the operation.
	BeforeRespondStream(ctx *RespondContext, rw http.ResponseWriter, request *http.Request, body io.Reader) *RespondContext
	RespondHook(ctx *RespondHookContext)
	RespondErrorHookContext(ctx *RespondErrorHookContext)
	RespondStreamHookContext(ctx *RespondStreamHookContext)
}

type (
	BeforeRespondFunc       = func(ctx *RespondContext, rw http.ResponseWriter, request *http.Request, body any) *RespondContext
	BeforeRespondErrorFunc  = func(ctx *RespondContext, rw http.ResponseWriter, request *http.Request, payload error) *RespondContext
	BeforeRespondStreamFunc = func(ctx *RespondContext, rw http.ResponseWriter, request *http.Request, body io.Reader) *RespondContext
	ResponseHookFunc        = func(ctx *RespondHookContext)
	ResponseErrorHookFunc   = func(ctx *RespondErrorHookContext)
	ResponseStreamHookFunc  = func(ctx *RespondStreamHookContext)
)

func (r *Responder) RegisterHook(hook RespondHook) {
	r.hooks = append(r.hooks, hook)
}

var _ RespondHook = (*respondHook)(nil)

type respondHook struct {
	readRequestLimit    int
	readRespondLimit    int
	filterRequest       FilterRequest
	filterRespondStream FilterRespond
	beforeRespond       BeforeRespondFunc
	beforeRespondError  BeforeRespondErrorFunc
	beforeRespondStream BeforeRespondStreamFunc
	onRespond           ResponseHookFunc
	onRespondError      ResponseErrorHookFunc
	onRespondStream     ResponseStreamHookFunc
}

func NewRespondHook(opts ...RespondHookOption) RespondHook {
	r := &respondHook{}
	for _, opt := range opts {
		opt.apply(r)
	}
	return r
}

func (r2 respondHook) AcceptRequestBodySize(r *http.Request) int {
	if r2.filterRequest != nil && r2.filterRequest(r) {
		return r2.readRequestLimit
	}
	return 0
}

func (r2 respondHook) AcceptResponseBodyStreamSize(contentType string, request *http.Request) int {
	if r2.filterRequest != nil && r2.filterRespondStream(contentType, request) {
		return r2.readRespondLimit
	}
	return 0
}

func (r2 respondHook) BeforeRespond(ctx *RespondContext, rw http.ResponseWriter, request *http.Request, body any) *RespondContext {
	if r2.beforeRespond != nil {
		return r2.beforeRespond(ctx, rw, request, body)
	}
	return ctx
}

func (r2 respondHook) BeforeRespondError(ctx *RespondContext, rw http.ResponseWriter, request *http.Request, errPayload error) *RespondContext {
	if r2.beforeRespondError != nil {
		return r2.beforeRespondError(ctx, rw, request, errPayload)
	}
	return ctx
}

func (r2 respondHook) BeforeRespondStream(ctx *RespondContext, rw http.ResponseWriter, request *http.Request, reader io.Reader) *RespondContext {
	if r2.beforeRespondStream != nil {
		return r2.beforeRespondStream(ctx, rw, request, reader)
	}
	return ctx
}

func (r2 respondHook) RespondHook(ctx *RespondHookContext) {
	if r2.onRespond != nil {
		r2.onRespond(ctx)
	}
}

func (r2 respondHook) RespondErrorHookContext(ctx *RespondErrorHookContext) {
	if r2.onRespondError != nil {
		r2.onRespondError(ctx)
	}
}

func (r2 respondHook) RespondStreamHookContext(ctx *RespondStreamHookContext) {
	if r2.onRespondStream != nil {
		r2.onRespondStream(ctx)
	}
}
