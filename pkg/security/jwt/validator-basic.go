package jwt

import (
	"crypto/rsa"
	"encoding/json"
	"github.com/BRBussy/goback/pkg/security/claims"
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
		rsaKeyPair: rsaKeyPair,
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
	jsonClaims, err := jwtObject.Verify(&b.rsaKeyPair.PublicKey)
	if err != nil {
		log.Error().Err(err).Msg("jwt verification failure")
		return nil, NewErrJWTVerificationFailure(err)
	}

	// unmarshal claims
	var serializedClaims claims.SerializedClaims
	if err := json.Unmarshal(jsonClaims, &serializedClaims); err != nil {
		log.Warn().Err(err).Msg("could not unmarshal claims")
		return nil, NewErrJSONUnmarshalError(err)
	}

	// check that claims are not expired
	if serializedClaims.Claims.Expired() {
		return nil, NewErrJWTExpired()
	}

	return &ValidateResponse{
		Claims: serializedClaims.Claims,
	}, nil
}
