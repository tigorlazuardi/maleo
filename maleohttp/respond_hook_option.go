package maleohttp

import (
	"net/http"
)

type RespondHookOption interface {
	apply(*respondHook)
}

type (
	RespondHookOptionBuilder []RespondHookOption
	respondHookOptionFunc    func(*respondHook)
)

func (r respondHookOptionFunc) apply(hook *respondHook) {
	r(hook)
}

func (h RespondHookOptionBuilder) apply(hook *respondHook) {
	for _, v := range h {
		v.apply(hook)
	}
}

type FilterRequest = func(*http.Request) bool

type FilterRespond = func(respondContentType string, r *http.Request) bool

// ReadRequestBodyLimit limits the number of bytes body being cloned. Defaults to 1MB.
//
// Negative value will make the hook clones all the body.
//
// Body will not be read if FilterRequest returns false.
func (h RespondHookOptionBuilder) ReadRequestBodyLimit(limit int) RespondHookOptionBuilder {
	return append(h, respondHookOptionFunc(func(r *respondHook) {
		r.readRequestLimit = limit
	}))
}

// ReadRespondBodyStreamLimit limits the number of bytes of respond body being cloned. Defaults to 1MB.
//
// Negative value will make the hook clones all the body.
//
// Body will not be read if FilterResponds returns false.
func (h RespondHookOptionBuilder) ReadRespondBodyStreamLimit(limit int) RespondHookOptionBuilder {
	return append(h, respondHookOptionFunc(func(r *respondHook) {
		r.readRespondLimit = limit
	}))
}

// FilterRequest filter requests whose body are going to be cloned. Defaults to filter only human readable content type.
func (h RespondHookOptionBuilder) FilterRequest(filter FilterRequest) RespondHookOptionBuilder {
	return append(h, respondHookOptionFunc(func(r *respondHook) {
		r.filterRequest = filter
	}))
}

// FilterRespondStream filter http server responds to clone. Defaults to filter only human readable content type.
func (h RespondHookOptionBuilder) FilterRespondStream(filter FilterRespond) RespondHookOptionBuilder {
	return append(h, respondHookOptionFunc(func(r *respondHook) {
		r.filterRespondStream = filter
	}))
}

// BeforeRespond is called after populating the handlers and values but before doing any operation to the response writer.
//
// You must not write any status code or body to the response writer.
// It will be set by the library.
//
// You are allowed to set any extra header or cookie to the response writer.
//
// You can modify everything else. Return the ctx to continue the operation.
func (h RespondHookOptionBuilder) BeforeRespond(before BeforeRespondFunc) RespondHookOptionBuilder {
	return append(h, respondHookOptionFunc(func(r *respondHook) {
		r.beforeRespond = before
	}))
}

// BeforeRespondError is called after populating the handlers and values but before doing any operation to the response writer.
//
// You must not write any status code or body to the response writer.
// It will be set by the library.
//
// You are allowed to set any extra header or cookie to the response writer.
//
// You can modify everything else. Return the ctx to continue the operation.
func (h RespondHookOptionBuilder) BeforeRespondError(before BeforeRespondErrorFunc) RespondHookOptionBuilder {
	return append(h, respondHookOptionFunc(func(r *respondHook) {
		r.beforeRespondError = before
	}))
}

// BeforeRespondStream is called after populating the handlers and values but before doing any operation to the response writer.
//
// You must not write any status code or body to the response writer.
// It will be set by the library.
//
// You are allowed to set any extra header or cookie to the response writer.
//
// You can modify everything else. Return the ctx to continue the operation.
func (h RespondHookOptionBuilder) BeforeRespondStream(before BeforeRespondStreamFunc) RespondHookOptionBuilder {
	return append(h, respondHookOptionFunc(func(r *respondHook) {
		r.beforeRespondStream = before
	}))
}

// OnRespond provides callback to be run after Responder writes the body.
//
// OnRespond callback is executed when maleohttp.Responder.Respond() is called.
//
// By default, the hook will use this api to call maleo to log the request and respond.
func (h RespondHookOptionBuilder) OnRespond(on ResponseHookFunc) RespondHookOptionBuilder {
	return append(h, respondHookOptionFunc(func(r *respondHook) {
		r.onRespond = on
	}))
}

// OnRespondError provides callback to be run after Responder writes the body.
//
// OnRespondError callback is executed when maleohttp.Responder.RespondError() is called.
//
// By default, the hook will use this api to call maleo to log the request and respond.
func (h RespondHookOptionBuilder) OnRespondError(on ResponseErrorHookFunc) RespondHookOptionBuilder {
	return append(h, respondHookOptionFunc(func(r *respondHook) {
		r.onRespondError = on
	}))
}

// OnRespondStream provides callback to be run after Responder writes the body.
//
// OnRespondStream callback is executed when maleohttp.Responder.RespondStream() is called.
//
// By default, the hook will use this api to call maleo to log the request and respond.
func (h RespondHookOptionBuilder) OnRespondStream(on ResponseStreamHookFunc) RespondHookOptionBuilder {
	return append(h, respondHookOptionFunc(func(r *respondHook) {
		r.onRespondStream = on
	}))
}
