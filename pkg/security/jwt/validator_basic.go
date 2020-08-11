package jwt

import (
	"github.com/BRBussy/goback/pkg/validate"
	"github.com/rs/zerolog/log"
)

type BasicValidator struct {
	requestValidator *validate.RequestValidator
}

func NewBasicValidator() *BasicValidator {
	return &BasicValidator{}
}

func (b *BasicValidator) Validate(request ValidateRequest) (*ValidateResponse, error) {
	if err := b.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &ValidateResponse{}, nil
}
