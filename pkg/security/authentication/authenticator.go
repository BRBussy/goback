package authentication

type Authenticator interface {
	Login(LoginRequest) (*LoginResponse, error)
}

type LoginRequest struct {
	EmailAddress string `validate:"required"`
	Password     string `validate:"required"`
}

type LoginResponse struct {
	JWT string
}
