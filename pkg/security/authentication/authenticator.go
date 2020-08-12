package authentication

import "github.com/BRBussy/goback/pkg/security/claims"

type Authenticator interface {
	Login(LoginRequest) (*LoginResponse, error)
	ValidateJWT(ValidateJWTRequest) (*ValidateJWTResponse, error)
}

const AuthenticatorServiceProviderName = "Authenticator"

type LoginRequest struct {
	Email    string `validate:"required"`
	Password string `validate:"required"`
}

type LoginResponse struct {
	JWT string
}

type ValidateJWTRequest struct {
	JWT string `validate:"required"`
}

type ValidateJWTResponse struct {
	Claims claims.Claims
}
