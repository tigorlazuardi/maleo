package maleohttp

import (
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/tigorlazuardi/maleo"
)

type RoundTripContext struct {
	Request     *http.Request
	RequestBody ClonedBody
	// Response is nil if there is an error to get HTTP Response.
	Response     *http.Response
	ResponseBody ClonedBody
	Error        error
	Caller       maleo.Caller
	Maleo        *maleo.Maleo
}

type RoundTrip struct {
	inner       http.RoundTripper
	hook        RoundTripHook
	maleo       *maleo.Maleo
	callerDepth int
}

const sep = string(os.PathSeparator)

var srcFile = strings.Join([]string{runtime.GOROOT(), sep, "src", sep, "net", sep, "http", sep, "client.go"}, "")

// RoundTrip implements http.RoundTripper interface.
func (rt *RoundTrip) RoundTrip(req *http.Request) (*http.Response, error) {
	var reqBody ClonedBody = NoopCloneBody{}
	req = req.Clone(req.Context())
	wantReqBody := rt.hook.AcceptRequestBodySize(req)
	if wantReqBody != 0 {
		reqBodyClone := wrapBodyCloner(req.Body, wantReqBody)
		req.Body = reqBodyClone
		reqBody = reqBodyClone
	}
	res, err := rt.inner.RoundTrip(req)

	caller := maleo.GetCaller(rt.callerDepth)
	// detect client.Get(), client.Head(), client.Post(), client.PostForm() request.
	if caller.File() == srcFile {
		caller = maleo.GetCaller(rt.callerDepth + 1)
	}
	ctx := &RoundTripContext{
		Request:      req,
		RequestBody:  reqBody,
		Response:     res,
		ResponseBody: NoopCloneBody{},
		Error:        err,
		Maleo:        rt.maleo,
		Caller:       caller,
	}
	if res != nil {
		wantResBody := rt.hook.AcceptResponseBodySize(req, res)
		if wantResBody != 0 {
			resBodyClone := wrapBodyCloner(res.Body, wantResBody)
			ctx.ResponseBody = resBodyClone
			resBodyClone.onClose(func(error) {
				rt.hook.ExecuteHook(ctx)
			})
			res.Body = resBodyClone
		}
	} else {
		rt.hook.ExecuteHook(ctx)
	}

	return res, err
}

// NewRoundTrip creates a new instance of http.RoundTripper implementation that wraps around the default http.RoundTripper.
// It provides logging support using maleo engine.
//
// By default, RoundTrip uses maleo's Global Instance.
// You may override this with maleohttp.Option.RoundTrip().Maleo(*maleo.Maleo)
//
// Caller by default points to where Client.Do(*http.Request) is called, but this assumes you don't use other
// http.RoundTripper or using custom client. if the caller location is incorrect,
// You may override this with maleohttp.Option.RoundTrip().AddCallerDepth(int) or maleohttp.Option.RoundTrip().CallerDepth(int)
//
// For reference, the default caller depth is 6.
func NewRoundTrip(opts ...RoundTripOption) *RoundTrip {
	rt := &RoundTrip{inner: http.DefaultTransport, hook: NewRoundTripHook(), maleo: maleo.Global(), callerDepth: 6}
	for _, v := range opts {
		v.apply(rt)
	}
	return rt
}

// WrapRoundTripper wraps the given http.RoundTripper with maleohttp.RoundTrip implementation to support logging with
// maleo engine.
//
// By default, RoundTrip uses maleo's Global Instance.
// You may override this with maleohttp.Option.RoundTrip().Maleo(*maleo.Maleo)
//
// Caller by default points to where Client.Do(*http.Request) is called, but this assumes you don't use other
// http.RoundTripper or using custom client. if the caller location is incorrect,
// You may override this with maleohttp.Option.RoundTrip().AddCallerDepth(int) or maleohttp.Option.RoundTrip().CallerDepth(int)
//
// For reference, the default caller depth is 6.
func WrapRoundTripper(rt http.RoundTripper, opts ...RoundTripOption) *RoundTrip {
	roundtrip := &RoundTrip{inner: rt, hook: NewRoundTripHook(), maleo: maleo.Global(), callerDepth: 6}
	for _, v := range opts {
		v.apply(roundtrip)
	}
	return roundtrip
}

// WrapHTTPClient wraps the given http.Client with maleohttp.RoundTrip implementation to support logging with
// maleo engine.
//
// if the passed client is nil, maleohttp will use http.DefaultClient.
// if the client.Transport is nil, maleohttp will use http.DefaultTransport as base.
//
// By default, RoundTrip uses maleo's Global Instance.
// You may override this with maleohttp.Option.RoundTrip().Maleo(*maleo.Maleo)
//
// Caller by default points to where Client.Do(*http.Request) is called, but this assumes you don't use other
// http.RoundTripper or not using custom client. if the caller location is incorrect,
// You may override this with maleohttp.Option.RoundTrip().AddCallerDepth(int) or maleohttp.Option.RoundTrip().CallerDepth(int)
//
// For reference, the default caller depth is 6.
func WrapHTTPClient(client *http.Client, opts ...RoundTripOption) *http.Client {
	if client == nil {
		client = http.DefaultClient
	}
	transport := client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	client.Transport = WrapRoundTripper(transport, opts...)
	return client
}

type RoundTripOption interface {
	apply(*RoundTrip)
}

type (
	RoundTripOptionFunc    func(*RoundTrip)
	RoundTripOptionBuilder []RoundTripOption
)

func (rto RoundTripOptionFunc) apply(rt *RoundTrip) {
	rto(rt)
}

func (rtob RoundTripOptionBuilder) apply(rt *RoundTrip) {
	for _, v := range rtob {
		v.apply(rt)
	}
}

// Maleo sets the maleo instance to use for logging.
func (rtob RoundTripOptionBuilder) Maleo(t *maleo.Maleo) RoundTripOptionBuilder {
	return append(rtob, RoundTripOptionFunc(func(rt *RoundTrip) {
		rt.maleo = t
	}))
}

// CallerDepth sets the caller depth to the caller stack.
func (rtob RoundTripOptionBuilder) CallerDepth(depth int) RoundTripOptionBuilder {
	return append(rtob, RoundTripOptionFunc(func(rt *RoundTrip) {
		rt.callerDepth = depth
	}))
}

// AddCallerDepth adds the caller depth to the caller stack.
func (rtob RoundTripOptionBuilder) AddCallerDepth(depth int) RoundTripOptionBuilder {
	return append(rtob, RoundTripOptionFunc(func(rt *RoundTrip) {
		rt.callerDepth += depth
	}))
}
