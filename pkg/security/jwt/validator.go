package jwt

type Validator interface {
	Validate(ValidateRequest) (*ValidateResponse, error)
}

type ValidateRequest struct {
	JWT string `validate:"required"`
}

type ValidateResponse struct {
}
