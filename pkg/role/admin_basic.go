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

func (b BasicAdmin) AddNewRole(request AddNewRoleRequest) (*AddNewRoleResponse, error) {
	// validate service request
	if err := b.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &AddNewRoleResponse{}, nil
}

func (b BasicAdmin) UpdateRole(request UpdateRoleRequest) (*UpdateRoleResponse, error) {
	// validate service request
	if err := b.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &UpdateRoleResponse{}, nil
}
