package authentication

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
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

// Talking about context:
// https://blog.questionable.services/article/map-string-interface/

func (a *Middleware) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// look for authentication header on request
		authenticationHeader := r.Header.Get("authentication")
		if authenticationHeader == "" {
			log.Warn().Msg("authentication header not set on request")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// try and validate the contents of the authentication header
		validateResponse, err := a.authenticator.ValidateJWT(
			ValidateJWTRequest{
				JWT: authenticationHeader,
			},
		)
		if err != nil {
			log.Warn().Err(err).Msg("jwt validation error")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// marshall claims and put into context
		marshalledClaims, err := json.Marshal(validateResponse.Claims)
		if err != nil {
			log.Warn().Err(err).Msg("error json marshalling claims")
			http.Error(w, "Unauthorized", http.StatusInternalServerError)
			return
		}

		// forward http request with claims context
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "Claims", marshalledClaims)))
	})
}
