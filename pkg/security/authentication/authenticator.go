package authentication

type Authenticator interface {
	Login(LoginRequest) (*LoginResponse, error)
}

const AuthenticatorServiceProviderName = "Authenticator"

type LoginRequest struct {
	Email    string `validate:"required"`
	Password string `validate:"required"`
}

type LoginResponse struct {
	JWT string
}
