package authorisation

import (
	"github.com/BRBussy/goback/pkg/role"
	"github.com/BRBussy/goback/pkg/user"
	"github.com/BRBussy/goback/pkg/validate"
	"github.com/rs/zerolog/log"
)

type BasicAuthorizer struct {
	requestValidator *validate.RequestValidator
	userStore        user.Store
	roleStore        role.Store
}

func NewBasicAuthorizer(
	userStore user.Store,
	roleStore role.Store,
) *BasicAuthorizer {
	return &BasicAuthorizer{
		requestValidator: validate.NewRequestValidator(),
		userStore:        userStore,
		roleStore:        roleStore,
	}
}

func (b *BasicAuthorizer) ConfirmServiceAccess(request ConfirmServiceAccessRequest) (*ConfirmServiceAccessResponse, error) {
	if err := b.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &ConfirmServiceAccessResponse{}, nil
}
