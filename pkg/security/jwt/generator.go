package jwt

import "github.com/BRBussy/goback/pkg/security/claims"

type Generator interface {
	Generate(GenerateRequest) (*GenerateResponse, error)
}

type GenerateRequest struct {
	Claims claims.Claims `validate:"required"`
}

type GenerateResponse struct {
	JWT string
}
