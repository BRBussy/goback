package user

import (
	"github.com/BRBussy/goback/pkg/role"
	"github.com/BRBussy/goback/pkg/validate"
	"github.com/rs/zerolog/log"
)

type BasicAdmin struct {
	validator *validate.RequestValidator
	userStore Store
	roleStore role.Store
}

func NewBasicAdmin(
	userStore Store,
	roleStore role.Store,
) *BasicAdmin {
	return &BasicAdmin{
		validator: validate.NewRequestValidator(),
		userStore: userStore,
	}
}

func (b BasicAdmin) AddNewUser(request AddNewUserRequest) (*AddNewUserResponse, error) {
	// validate service request
	if err := b.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &AddNewUserResponse{}, nil
}
