package authentication

import (
	"net/http"
)

type Middleware struct {
	authenticator Authenticator
}

func NewMiddleware(
	authenticator Authenticator,
) *Middleware {
	return &Middleware{
		authenticator: authenticator,
	}
}

func (a *Middleware) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
