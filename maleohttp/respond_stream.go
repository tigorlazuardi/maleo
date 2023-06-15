package maleohttp

import (
	"io"
	"net/http"

	"github.com/tigorlazuardi/maleo"
)

// RespondStream writes the given stream to the http.ResponseWriter.
//
// If the stream implements maleo.HTTPCodeHint, the status code will be set to the value returned by the maleo.HTTPCodeHint.
//
// There's a special case if you pass http.NoBody as body, there will be no respond body related operations executed.
// StatusCode default value is set to http.StatusNoContent. You can still override this output by setting the
// related RespondOption.
// With http.NoBody as body, Maleohttp will immediately respond with status code after RespondOption are evaluated
// and end the process.
//
// Body of nil will be treated as http.NoBody.
func (r *Responder) RespondStream(rw http.ResponseWriter, request *http.Request, contentType string, body io.Reader, opts ...RespondOption) {
	var (
		statusCode = http.StatusOK
		err        error
	)
	if body == nil {
		body = http.NoBody
	}
	if ch, ok := body.(maleo.HTTPCodeHint); ok {
		statusCode = ch.HTTPCode()
	}
	opt := r.buildOptionStream(statusCode, rw, request, body, opts...)
	if len(r.hooks) > 0 {
		var clone ClonedBody = NoopCloneBody{}
		count := r.hooks.CountMaximumRespondBodyRead(contentType, request)
		if count != 0 {
			s := wrapBodyCloner(body, count)
			if body != http.NoBody {
				body = s
			}
			clone = s
		}
		defer func() {
			var requestBody ClonedBody = NoopCloneBody{}
			if b, ok := request.Body.(ClonedBody); ok {
				requestBody = b
			} else if c := clonedBodyFromContext(request.Context()); c != nil {
				requestBody = c
			}
			hookContext := &RespondStreamHookContext{
				baseHook: &baseHook{
					Context:        opt,
					Request:        request,
					RequestBody:    requestBody,
					ResponseStatus: opt.StatusCode,
					ResponseHeader: rw.Header(),
					Maleo:          r.maleo,
					Error:          err,
				},
				ResponseBody: RespondStreamBody{
					Value:       clone,
					ContentType: contentType,
				},
			}
			for _, hook := range r.hooks {
				hook.RespondStreamHookContext(hookContext)
			}
		}()
	}
	if body == http.NoBody {
		rw.WriteHeader(opt.StatusCode)
		return
	}

	compressed, ok := opt.StreamCompressor.StreamCompress(contentType, body)
	if ok {
		rw.Header().Set("Content-Encoding", opt.StreamCompressor.ContentEncoding())
		body = compressed
	}
	if len(contentType) > 0 {
		rw.Header().Set("Content-Type", contentType)
	}
	rw.WriteHeader(opt.StatusCode)
	_, err = io.Copy(rw, body)
}
