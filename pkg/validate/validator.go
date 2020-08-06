package validate

import (
	"fmt"
	goValidator "gopkg.in/go-playground/validator.v9"
)

type RequestValidator struct {
	validate *goValidator.Validate
}

func NewRequestValidator() *RequestValidator {
	return &RequestValidator{
		validate: goValidator.New(),
	}
}

func (v *RequestValidator) Validate(request interface{}) error {
	var reasons []string

	if err := v.validate.Struct(request); err != nil {
		validationErrors := err.(goValidator.ValidationErrors)
		for key, value := range validationErrors {
			reasons = append(reasons, fmt.Sprintf("'%d' %s", key, value))
		}
	}

	if len(reasons) > 0 {
		return NewErrRequestNotValid(reasons)
	}
	return nil
}
