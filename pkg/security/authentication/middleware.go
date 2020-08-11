package authentication

import (
	"github.com/BRBussy/goback/pkg/security/jwt"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Middleware struct {
	authenticator Authenticator
	jwtValidator  jwt.Validator
}

func NewMiddleware(
	authenticator Authenticator,
	jwtValidator jwt.Validator,
) *Middleware {
	return &Middleware{
		authenticator: authenticator,
		jwtValidator:  jwtValidator,
	}
}

func (a *Middleware) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// look for authentication header on request
		authenticationHeader := r.Header.Get("authentication")
		if authenticationHeader == "" {
			log.Warn().Msg("authentication header not set on request")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// try and validate the contents of the authentication header as a jwt
		if _, err := a.jwtValidator.Validate(
			jwt.ValidateRequest{
				JWT: authenticationHeader,
			},
		); err != nil {
			log.Warn().Err(err).Msg("error validating authentication header jwt")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
