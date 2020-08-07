package role

import (
	"github.com/BRBussy/goback/pkg/validate"
	"github.com/rs/zerolog/log"
)

type BasicAdmin struct {
	requestValidator *validate.RequestValidator
	roleStore        Store
}

func NewBasicAdmin(
	roleStore Store,
) *BasicAdmin {
	return &BasicAdmin{
		requestValidator: validate.NewRequestValidator(),
		roleStore:        roleStore,
	}
}

func (b BasicAdmin) AddNewRole(request AddNewRoleRequest) (*AddNewRoleResponse, error) {
	if err := b.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &AddNewRoleResponse{}, nil
}

func (b BasicAdmin) UpdateRole(request UpdateRoleRequest) (*UpdateRoleResponse, error) {
	if err := b.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &UpdateRoleResponse{}, nil
}
