package maleohttp

import (
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

	BeforeRespond(ctx *RespondContext, request *http.Request) *RespondContext
	RespondHook(ctx *RespondHookContext)
	RespondErrorHookContext(ctx *RespondErrorHookContext)
	RespondStreamHookContext(ctx *RespondStreamHookContext)
}

type (
	BeforeRespondFunc      = func(ctx *RespondContext, request *http.Request) *RespondContext
	ResponseHookFunc       = func(ctx *RespondHookContext)
	ResponseErrorHookFunc  = func(ctx *RespondErrorHookContext)
	ResponseStreamHookFunc = func(ctx *RespondStreamHookContext)
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

func (r2 respondHook) BeforeRespond(ctx *RespondContext, request *http.Request) *RespondContext {
	if r2.beforeRespond == nil {
		return ctx
	}
	return r2.beforeRespond(ctx, request)
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
