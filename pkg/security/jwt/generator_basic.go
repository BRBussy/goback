package jwt

import (
	"crypto/rsa"
	"encoding/json"
	"github.com/BRBussy/goback/pkg/exception"
	"github.com/BRBussy/goback/pkg/security/claims"
	"github.com/BRBussy/goback/pkg/validate"
	"github.com/rs/zerolog/log"
	"gopkg.in/square/go-jose.v2"
)

type BasicGenerator struct {
	requestValidator *validate.RequestValidator
	tokenSigner      jose.Signer
}

func NewBasicGenerator(
	rsaKeyPair *rsa.PrivateKey,
) *BasicGenerator {
	// create a new signer with the given private key.
	joseSigner, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.PS512,
			Key:       rsaKeyPair,
		},
		nil,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("error generating new jose signer")
	}

	return &BasicGenerator{
		requestValidator: validate.NewRequestValidator(),
		tokenSigner:      joseSigner,
	}
}

func (b *BasicGenerator) Generate(request GenerateRequest) (*GenerateResponse, error) {
	if err := b.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// marshal the given claims to generate jwt payload
	jwtPayload, err := json.Marshal(
		claims.SerializedClaims{
			Claims: request.Claims,
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("error marshalling claims")
		return nil, NewErrJSONMarshalError(err)
	}

	// sign the jwt payload
	signedJWTObject, err := b.tokenSigner.Sign(jwtPayload)
	if err != nil {
		log.Error().Err(err).Msg("could not sign payload")
		return nil, NewErrSigningError(err)
	}

	// serialize the signed JWT Object
	signedJWT, err := signedJWTObject.CompactSerialize()
	if err != nil {
		log.Error().Err(err).Msg("error serializing signed JWT object")
		return nil, exception.NewErrUnexpected(err)
	}

	return &GenerateResponse{JWT: signedJWT}, nil
}
