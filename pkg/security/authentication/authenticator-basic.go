package authentication

import (
	"encoding/json"
	"errors"
	"github.com/BRBussy/goback/pkg/exception"
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
	jwtValidator     jwt.Validator
}

func NewBasicAuthenticator(
	userStore user.Store,
	jwtGenerator jwt.Generator,
	jwtValidator jwt.Validator,
) *BasicAuthenticator {
	return &BasicAuthenticator{
		requestValidator: validate.NewRequestValidator(),
		userStore:        userStore,
		jwtGenerator:     jwtGenerator,
		jwtValidator:     jwtValidator,
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
			Filter: filter.NewTextExactFilter(
				"email",
				request.Email,
			),
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
			Claims: claims.NewClaims(
				retrieveUserResponse.User.ID,
				time.Now().Add(24*time.Hour),
			),
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

func (b BasicAuthenticator) ValidateJWT(request ValidateJWTRequest) (*ValidateJWTResponse, error) {
	if err := b.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	validateResponse, err := b.jwtValidator.Validate(
		jwt.ValidateRequest{JWT: request.JWT},
	)
	if errors.Is(err, &jwt.ErrJWTInvalid{}) ||
		errors.Is(err, &jwt.ErrJWTVerificationFailure{}) {
		return nil, NewErrJWTInvalid()
	} else if err != nil {
		log.Error().Err(err).Msg("error validating jwt")
		return nil, exception.NewErrUnexpected(err)
	}

	// unmarshal claims
	var userClaims claims.Claims
	if err := json.Unmarshal(validateResponse.JSONPayload, &userClaims); err != nil {
		log.Warn().Err(err).Msg("could not unmarshal jwt json payload to claims")
		return nil, NewErrJSONUnmarshalError(err)
	}

	// check that claims are not expired
	if userClaims.Expired() {
		return nil, NewErrJWTExpired()
	}

	return &ValidateJWTResponse{
		Claims: userClaims,
	}, nil
}
