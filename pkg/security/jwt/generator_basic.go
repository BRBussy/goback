package jwt

import (
	"github.com/BRBussy/goback/pkg/validate"
	"github.com/rs/zerolog/log"
)

type BasicGenerator struct {
	requestValidator *validate.RequestValidator
}

func NewBasicGenerator(
	requestValidator *validate.RequestValidator,
) *BasicGenerator {
	return &BasicGenerator{
		requestValidator: requestValidator,
	}
}

func (b *BasicGenerator) Generate(request GenerateRequest) (*GenerateResponse, error) {
	if err := b.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &GenerateResponse{}, nil
}
