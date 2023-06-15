package maleohttp

import (
	"io"
	"net/http"
	"strconv"

	"github.com/tigorlazuardi/maleo"
)

type errString string

func (e errString) Error() string {
	return string(e)
}

const errInternalServerError errString = "Internal Server Error"

// RespondError writes the given error to the http.ResponseWriter.
//
// errPayload is expected to be a serializable type.
//
// HTTP Status code by default is http.StatusInternalServerError. If error implements maleo.HTTPCodeHint, the status code will be set to the
// value returned by the maleo.HTTPCodeHint method. If the maleohttp.Option.StatusCode RespondOption is set, it will override
// the status regardless of the maleo.HTTPCodeHint.
//
// if err is nil, it will be replaced with "Internal Server Error" message. It is done this way, because the library
// assumes that you mishandled the method and to prevent sending empty values, a generic Internal Server Error message
// will be sent instead. If you wish to send an empty response, use Respond with http.NoBody as body.
func (r Responder) RespondError(rw http.ResponseWriter, request *http.Request, errPayload error, opts ...RespondOption) {
	var (
		ctx            = request.Context()
		encodedBody    []byte
		err            error
		statusCode     = maleo.Query.GetHTTPCode(errPayload)
		compressedBody []byte
	)
	if errPayload == nil {
		errPayload = errInternalServerError
	}
	opt := r.buildOptionError(statusCode, rw, request, errPayload, opts...)
	if len(r.hooks) > 0 {
		defer func() {
			var requestBody ClonedBody = NoopCloneBody{}
			if b, ok := request.Body.(ClonedBody); ok {
				requestBody = b
			} else if c := clonedBodyFromContext(request.Context()); c != nil {
				requestBody = c
			}
			hookContext := &RespondErrorHookContext{
				baseHook: &baseHook{
					Context:        opt,
					Request:        request,
					RequestBody:    requestBody,
					ResponseStatus: opt.StatusCode,
					ResponseHeader: rw.Header(),
					Maleo:          r.maleo,
					Error:          err,
				},
				ResponseBody: RespondErrorBody{
					PreEncoded:     errPayload,
					PostEncoded:    encodedBody,
					PostCompressed: compressedBody,
				},
			}
			for _, hook := range r.hooks {
				hook.RespondErrorHookContext(hookContext)
			}
		}()
	}
	body := r.errorTransformer.ErrorBodyTransform(ctx, errPayload)
	if body == nil {
		rw.WriteHeader(opt.StatusCode)
		return
	}
	encodedBody, err = opt.Encoder.Encode(body)
	if err != nil {
		const errMsg = "ENCODING ERROR"
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		rw.WriteHeader(http.StatusInternalServerError)
		_, err = io.WriteString(rw, errMsg)
		encodedBody = []byte(errMsg)
		return
	}
	contentType := opt.Encoder.ContentType()
	if contentType != "" {
		rw.Header().Set("Content-Type", contentType)
	}
	compressedBody, ok, err := opt.Compressor.Compress(encodedBody)
	if err != nil {
		_ = r.maleo.Wrap(err).Caller(opt.Caller).Level(maleo.WarnLevel).Log(ctx)
		rw.Header().Set("Content-Length", strconv.Itoa(len(encodedBody)))
		rw.WriteHeader(opt.StatusCode)
		_, err = rw.Write(encodedBody)
		return
	}
	if ok {
		contentEncoding := opt.Compressor.ContentEncoding()
		rw.Header().Set("Content-Encoding", contentEncoding)
		rw.Header().Set("Content-Length", strconv.Itoa(len(compressedBody)))
		rw.WriteHeader(opt.StatusCode)
		_, err = rw.Write(compressedBody)
		return
	}
	rw.Header().Set("Content-Length", strconv.Itoa(len(encodedBody)))
	rw.WriteHeader(opt.StatusCode)
	_, err = rw.Write(encodedBody)
}
