package authorisation

import (
	"github.com/BRBussy/goback/pkg/role"
	"github.com/BRBussy/goback/pkg/user"
	"github.com/BRBussy/goback/pkg/validate"
	"github.com/rs/zerolog/log"
)

type BasicAuthoriser struct {
	requestValidator *validate.RequestValidator
	userStore        user.Store
	roleStore        role.Store
}

func NewBasicAuthoriser(
	userStore user.Store,
	roleStore role.Store,
) *BasicAuthoriser {
	return &BasicAuthoriser{
		requestValidator: validate.NewRequestValidator(),
		userStore:        userStore,
		roleStore:        roleStore,
	}
}

func (b *BasicAuthoriser) ConfirmServiceAccess(request ConfirmServiceAccessRequest) (*ConfirmServiceAccessResponse, error) {
	if err := b.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &ConfirmServiceAccessResponse{}, nil
}
