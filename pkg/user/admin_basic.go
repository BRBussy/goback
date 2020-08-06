package user

import (
	"github.com/BRBussy/goback/pkg/validate"
	"github.com/rs/zerolog/log"
)

type BasicAdmin struct {
	validator *validate.RequestValidator
	userStore Store
}

func NewBasicAdmin(
	userStore Store,
) *BasicAdmin {
	return &BasicAdmin{
		validator: validate.NewRequestValidator(),
		userStore: userStore,
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
