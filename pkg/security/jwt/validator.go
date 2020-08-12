package jwt

import (
	"encoding/json"
)

type Validator interface {
	Validate(ValidateRequest) (*ValidateResponse, error)
}

type ValidateRequest struct {
	JWT string `validate:"required"`
}

type ValidateResponse struct {
	JSONPayload json.RawMessage
}
