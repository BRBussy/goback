package authorisation

import "github.com/BRBussy/goback/pkg/security/claims"

type Authorizer interface {
	ConfirmServiceAccess(ConfirmServiceAccessRequest) (*ConfirmServiceAccessResponse, error)
}

type ConfirmServiceAccessRequest struct {
	Claims  claims.Claims `validate:"required"`
	Service string        `validate:"required"`
}

type ConfirmServiceAccessResponse struct {
	Result bool
}
