package jwt

import (
	"crypto/rsa"
	"github.com/BRBussy/goback/pkg/validate"
	"github.com/rs/zerolog/log"
	"gopkg.in/square/go-jose.v2"
)

type BasicValidator struct {
	requestValidator *validate.RequestValidator
	rsaKeyPair       *rsa.PrivateKey
}

func NewBasicValidator(
	rsaKeyPair *rsa.PrivateKey,
) *BasicValidator {
	return &BasicValidator{
		requestValidator: validate.NewRequestValidator(),
		rsaKeyPair:       rsaKeyPair,
	}
}

func (b *BasicValidator) Validate(request ValidateRequest) (*ValidateResponse, error) {
	if err := b.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// Parse the given jwt - success means the given jwt string is a jwt
	jwtObject, err := jose.ParseSigned(request.JWT)
	if err != nil {
		log.Error().Err(err).Msg("jwt parsing failure")
		return nil, NewErrJWTInvalid(err)
	}

	// Verify jwt signature and retrieve payload (i.e. json marshalled claims)
	// Failure indicates jwt was damaged or tampered with
	jsonPayload, err := jwtObject.Verify(&b.rsaKeyPair.PublicKey)
	if err != nil {
		log.Error().Err(err).Msg("jwt verification failure")
		return nil, NewErrJWTVerificationFailure()
	}

	return &ValidateResponse{
		JSONPayload: jsonPayload,
	}, nil
}
