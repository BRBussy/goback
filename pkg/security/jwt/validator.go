package jwt

import "github.com/BRBussy/goback/pkg/security/claims"

type Validator interface {
	Validate(ValidateRequest) (*ValidateResponse, error)
}

type ValidateRequest struct {
	JWT string `validate:"required"`
}

type ValidateResponse struct {
	Claims claims.Claims
}
