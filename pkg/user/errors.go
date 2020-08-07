package user

import "strings"

type ErrEmailAddressAlreadyUsed struct{}

func NewErrEmailAddressAlreadyUsed() *ErrEmailAddressAlreadyUsed {
	return &ErrEmailAddressAlreadyUsed{}
}

func (e *ErrEmailAddressAlreadyUsed) Error() string {
	return "email address already used"
}

type ErrUserDoesNotExist struct{}

func NewErrUserDoesNotExist() *ErrUserDoesNotExist {
	return &ErrUserDoesNotExist{}
}

func (e *ErrUserDoesNotExist) Error() string {
	return "user does not exist"
}

type ErrUserNotValid struct {
	ReasonsInvalid []string
}

func NewErrUserNotValid(
	reasonsInvalid []string,
) *ErrUserNotValid {
	return &ErrUserNotValid{
		ReasonsInvalid: reasonsInvalid,
	}
}

func (e *ErrUserNotValid) Error() string {
	return "user not valid: " + strings.Join(e.ReasonsInvalid, ", ")
}
