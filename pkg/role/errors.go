package role

import "strings"

type ErrRoleAlreadyExists struct{}

func NewErrRoleAlreadyExists() *ErrRoleAlreadyExists {
	return &ErrRoleAlreadyExists{}
}

func (e *ErrRoleAlreadyExists) Error() string {
	return "role already exists"
}

type ErrRoleDoesNotExist struct{}

func NewErrRoleDoesNotExist() *ErrRoleDoesNotExist {
	return &ErrRoleDoesNotExist{}
}

func (e *ErrRoleDoesNotExist) Error() string {
	return "role does not exist"
}

type ErrRoleNotValid struct {
	ReasonsInvalid []string
}

func NewErrRoleNotValid(
	reasonsInvalid []string,
) *ErrRoleNotValid {
	return &ErrRoleNotValid{
		ReasonsInvalid: reasonsInvalid,
	}
}

func (e *ErrRoleNotValid) Error() string {
	return "role not valid: " + strings.Join(e.ReasonsInvalid, ", ")
}
