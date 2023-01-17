package maleohttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/tigorlazuardi/maleo"
)

func NewLoggerHook(opts ...RespondHookOption) RespondHook {
	return NewRespondHook(append(defaultLoggerOptions(), opts...)...)
}

func defaultLoggerOptions() RespondHookOptionBuilder {
	return Option.RespondHook().
		FilterRequest(func(r *http.Request) bool {
			return isHumanReadable(r.Header.Get("Content-Type"))
		}).
		ReadRequestBodyLimit(1024 * 1024).
		ReadRespondBodyStreamLimit(1024 * 1024).
		FilterRespondStream(func(respondContentType string, r *http.Request) bool {
			return isHumanReadable(respondContentType)
		}).
		OnRespond(defaultLoggerRespond).
		OnRespondError(defaultLoggerRespondError).
		OnRespondStream(defaultLoggerRespondStream)
}

func defaultLoggerRespond(ctx *RespondHookContext) {
	fields := buildLoggerFields(ctx.baseHook, ctx.ResponseBody.PostEncoded, false)
	message := fmt.Sprintf("%s %s %s", ctx.Request.Method, ctx.Request.URL.String(), ctx.Request.Proto)
	if ctx.Error != nil {
		_ = ctx.Maleo.Wrap(ctx.Error).Level(maleo.ErrorLevel).Code(ctx.ResponseStatus).Message(message).Caller(ctx.Context.Caller).Context(fields).Log(ctx.Request.Context())
		return
	}
	ctx.Maleo.NewEntry(message).Level(maleo.InfoLevel).Code(ctx.ResponseStatus).Caller(ctx.Context.Caller).Context(fields).Log(ctx.Request.Context())
}

func defaultLoggerRespondError(ctx *RespondErrorHookContext) {
	fields := buildLoggerFields(ctx.baseHook, ctx.ResponseBody.PostEncoded, false)
	message := fmt.Sprintf("%s %s %s", ctx.Request.Method, ctx.Request.URL.String(), ctx.Request.Proto)
	if ctx.Error != nil {
		_ = ctx.Maleo.Wrap(ctx.Error).Code(ctx.ResponseStatus).Message(message).Caller(ctx.Context.Caller).Context(fields).Log(ctx.Request.Context())
		return
	}
	_ = ctx.Maleo.Wrap(ctx.ResponseBody.PreEncoded).Code(ctx.ResponseStatus).Caller(ctx.Context.Caller).Context(fields).Log(ctx.Request.Context())
}

func defaultLoggerRespondStream(ctx *RespondStreamHookContext) {
	fields := buildLoggerFields(ctx.baseHook, ctx.ResponseBody.Value.CloneBytes(), ctx.ResponseBody.Value.Truncated())
	message := fmt.Sprintf("%s %s %s", ctx.Request.Method, ctx.Request.URL.String(), ctx.Request.Proto)
	if ctx.Error != nil {
		_ = ctx.Maleo.Wrap(ctx.Error).Level(maleo.ErrorLevel).Code(ctx.ResponseStatus).Message(message).Caller(ctx.Context.Caller).Context(fields).Log(ctx.Request.Context())
		return
	}
	ctx.Maleo.NewEntry(message).Level(maleo.InfoLevel).Code(ctx.ResponseStatus).Caller(ctx.Context.Caller).Context(fields).Log(ctx.Request.Context())
}

func buildLoggerFields(hook *baseHook, respBody []byte, truncated bool) maleo.F {
	url := hook.Request.Host + hook.Request.URL.String()
	requestFields := maleo.F{
		"method": hook.Request.Method,
		"url":    url,
	}
	if len(hook.Request.Header) > 0 {
		requestFields["headers"] = hook.Request.Header
	}

	if hook.RequestBody.Len() > 0 {
		contentType := hook.Request.Header.Get("Content-Type")
		switch {
		case hook.RequestBody.Truncated():
			requestFields["body"] = fmt.Sprintf("%s (truncated)", hook.RequestBody.String())
		case strings.Contains(contentType, "application/json") && isJson(hook.RequestBody.Bytes()):
			requestFields["body"] = json.RawMessage(hook.RequestBody.CloneBytes())
		case contentType == "" && isJsonLite(hook.RequestBody.Bytes()) && isJson(hook.RequestBody.Bytes()):
			requestFields["body"] = json.RawMessage(hook.RequestBody.CloneBytes())
		default:
			requestFields["body"] = hook.RequestBody.String()
		}
	}

	responseFields := maleo.F{
		"status": hook.ResponseStatus,
	}
	if len(hook.ResponseHeader) > 0 {
		responseFields["headers"] = hook.ResponseHeader
	}
	if len(respBody) > 0 {
		contentType := hook.ResponseHeader.Get("Content-Type")
		switch {
		case truncated:
			responseFields["body"] = fmt.Sprintf("%s (truncated)", hook.RequestBody.String())
		case strings.Contains(contentType, "application/json") && isJson(respBody):
			responseFields["body"] = json.RawMessage(respBody)
		case contentType == "" && isJsonLite(respBody) && isJson(respBody):
			responseFields["body"] = json.RawMessage(respBody)
		default:
			responseFields["body"] = string(respBody)
		}
	}

	return maleo.F{
		"request":  requestFields,
		"response": responseFields,
	}
}

func isJsonLite(b []byte) bool {
	if len(b) < 2 {
		return false
	}
	return (b[0] == '{' || b[0] == '[') && (b[len(b)-1] == '}' || b[len(b)-1] == ']')
}
