package user

import (
	"github.com/BRBussy/goback/pkg/role"
	"github.com/BRBussy/goback/pkg/validate"
	"github.com/rs/zerolog/log"
)

type BasicAdmin struct {
	requestValidator *validate.RequestValidator
	userStore        Store
	roleStore        role.Store
}

func NewBasicAdmin(
	userStore Store,
	roleStore role.Store,
) *BasicAdmin {
	return &BasicAdmin{
		requestValidator: validate.NewRequestValidator(),
		userStore:        userStore,
		roleStore:        roleStore,
	}
}

func (b BasicAdmin) AddNewUser(request AddNewUserRequest) (*AddNewUserResponse, error) {
	if err := b.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &AddNewUserResponse{}, nil
}

func (b BasicAdmin) UpdateUser(request UpdateUserRequest) (*UpdateUserResponse, error) {
	if err := b.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &UpdateUserResponse{}, nil
}
