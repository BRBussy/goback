package role

import (
	"github.com/BRBussy/goback/pkg/validate"
	"github.com/rs/zerolog/log"
)

type BasicAdmin struct {
	validator *validate.RequestValidator
	roleStore Store
}

func NewBasicAdmin(
	roleStore Store,
) *BasicAdmin {
	return &BasicAdmin{
		validator: validate.NewRequestValidator(),
		roleStore: roleStore,
	}
}

func (b BasicAdmin) Get(request GetRequest) (*GetResponse, error) {
	// validate service request
	if err := b.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &GetResponse{}, nil
}

func (b BasicAdmin) Set(request SetRequest) (*SetResponse, error) {
	// validate service request
	if err := b.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &SetResponse{}, nil
}
