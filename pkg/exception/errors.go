package exception

import "fmt"

type ErrUnexpected struct {
	Err error
}

func NewErrUnexpected(err error) *ErrUnexpected {
	return &ErrUnexpected{Err: err}
}

func (e *ErrUnexpected) Error() string {
	return fmt.Sprintf("unexpected backend error: %v", e.Err)
}

func (e *ErrUnexpected) Unwrap() error {
	return e.Err
}

type ErrNoChangesMade struct {
}

func NewErrNoChangesMade() *ErrNoChangesMade {
	return &ErrNoChangesMade{}
}

func (e *ErrNoChangesMade) Error() string {
	return "no changes made"
}
