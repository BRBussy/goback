package authorisation

import (
	"net/http"
)

type Middleware struct {
	authoriser Authoriser
}

func NewMiddleware(
	authoriser Authoriser,
) *Middleware {
	return &Middleware{
		authoriser: authoriser,
	}
}

func (a *Middleware) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
