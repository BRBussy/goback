package authentication

import (
	"github.com/BRBussy/goback/pkg/mongo/filter"
	"github.com/BRBussy/goback/pkg/security/jwt"
	"github.com/BRBussy/goback/pkg/user"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Middleware struct {
	authenticator Authenticator
	jwtValidator  jwt.Validator
	userStore     user.Store
}

func NewMiddleware(
	authenticator Authenticator,
	jwtValidator jwt.Validator,
	userStore user.Store,
) *Middleware {
	return &Middleware{
		authenticator: authenticator,
		jwtValidator:  jwtValidator,
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

		// try and validate the contents of the authentication header as a jwt
		validateResponse, err := a.jwtValidator.Validate(
			jwt.ValidateRequest{
				JWT: authenticationHeader,
			},
		)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Warn().Err(err).Msg("unauthorized request")
			return
		}

		// use the returned claims to try and retrieve the requesting user
		retrieveResponse, err := a.userStore.Retrieve(
			user.RetrieveRequest{
				Filter: filter.NewTextExactFilter(
					"id",
					validateResponse.Claims.Expired(),
				),
			},
		)

		next.ServeHTTP(w, r)
	})
}
