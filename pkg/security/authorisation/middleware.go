package authorisation

import (
	"net/http"
)

type Middleware struct {
	authorizer Authorizer
}

func NewMiddleware(
	authorizer Authorizer,
) *Middleware {
	return &Middleware{
		authorizer: authorizer,
	}
}

func (a *Middleware) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
