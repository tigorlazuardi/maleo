package maleohttp

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func (r *Responder) RequestBodyCloner() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			size := r.hooks.CountMaximumRequestBodyRead(request)
			if size != 0 {
				cloner := wrapBodyCloner(request.Body, size)
				request.Body = cloner
				ctx := contextWithClonedBody(request.Context(), cloner)
				request = request.WithContext(ctx)
			}
			next.ServeHTTP(writer, request)
		})
	}
}
