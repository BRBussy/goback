package authentication

import (
	"github.com/BRBussy/goback/pkg/user"
	"github.com/BRBussy/goback/pkg/validate"
)

type BasicAuthenticator struct {
	userStore        user.Store
	requestValidator *validate.RequestValidator
}

func NewBasicAuthenticator(
	userStore user.Store,
) *BasicAuthenticator {
	return &BasicAuthenticator{
		requestValidator: validate.NewRequestValidator(),
		userStore:        userStore,
	}
}

func (b BasicAuthenticator) Login(request LoginRequest) (*LoginResponse, error) {
	panic("implement me")
}
