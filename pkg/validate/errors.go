package validate

import "strings"

type ErrRequestNotValid struct {
	Reasons []string
}

func NewErrRequestNotValid(reasons []string) *ErrRequestNotValid {
	return &ErrRequestNotValid{Reasons: reasons}
}

func (e *ErrRequestNotValid) Error() string {
	return "request not valid: " + strings.Join(e.Reasons, ", ")
}
