package authentication

import (
	"errors"
	"github.com/BRBussy/goback/pkg/mongo"
	"github.com/BRBussy/goback/pkg/mongo/filter"
	"github.com/BRBussy/goback/pkg/security/claims"
	"github.com/BRBussy/goback/pkg/security/jwt"
	"github.com/BRBussy/goback/pkg/user"
	"github.com/BRBussy/goback/pkg/validate"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type BasicAuthenticator struct {
	requestValidator *validate.RequestValidator
	userStore        user.Store
	jwtGenerator     jwt.Generator
}

func NewBasicAuthenticator(
	userStore user.Store,
	jwtGenerator jwt.Generator,
) *BasicAuthenticator {
	return &BasicAuthenticator{
		requestValidator: validate.NewRequestValidator(),
		userStore:        userStore,
		jwtGenerator:     jwtGenerator,
	}
}

func (b BasicAuthenticator) Login(request LoginRequest) (*LoginResponse, error) {
	if err := b.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// try and retrieve the user by the given email address
	retrieveUserResponse, err := b.userStore.Retrieve(
		user.RetrieveRequest{
			Filter: filter.NewEmailFilter(request.EmailAddress),
		},
	)
	if err != nil {
		if errors.Is(err, &mongo.ErrNotFound{}) {
			log.Warn().Msg("attempted login from non-existent user")
		} else {
			log.Error().Err(err).Msg("error retrieving user for login")
		}
		return nil, NewErrLoginFailed()
	}

	// check password is correct by comparing it with the stored hash
	if err := bcrypt.CompareHashAndPassword(
		retrieveUserResponse.User.Password,
		[]byte(request.Password),
	); err != nil {
		return nil, NewErrLoginFailed()
	}

	// generate a login claims JWT
	generateTokenResponse, err := b.jwtGenerator.Generate(
		jwt.GenerateRequest{
			Claims: claims.Login{
				UserID:         retrieveUserResponse.User.ID,
				ExpirationTime: time.Now().Add(24 * time.Hour),
			},
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to generate login claims jwt")
		return nil, NewErrLoginFailed()
	}

	return &LoginResponse{
		JWT: generateTokenResponse.JWT,
	}, nil
}
