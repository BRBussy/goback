package mongo

import (
	"fmt"
	"strings"
)

type ErrConnectionError struct {
	Err error
}

func NewErrConnectionError(err error) *ErrConnectionError {
	return &ErrConnectionError{Err: err}
}

func (e *ErrConnectionError) Error() string {
	return fmt.Sprintf("mongo connection error: %v", e.Err)
}

func (e *ErrConnectionError) Unwrap() error {
	return e.Err
}

type ErrCloseConnectionError struct {
	Err error
}

func NewErrCloseConnectionError(err error) *ErrCloseConnectionError {
	return &ErrCloseConnectionError{Err: err}
}

func (e *ErrCloseConnectionError) Error() string {
	return fmt.Sprintf("mongo close connection error: %v", e.Err)
}

func (e *ErrCloseConnectionError) Unwrap() error {
	return e.Err
}

type ErrPingError struct {
	Err error
}

func NewErrPingError(err error) *ErrConnectionError {
	return &ErrConnectionError{Err: err}
}

func (e *ErrPingError) Error() string {
	return fmt.Sprintf("mongo ping error: %v", e.Err)
}

func (e *ErrPingError) Unwrap() error {
	return e.Err
}

type ErrInvalidConfig struct {
	Reasons []string
}

func NewErrInvalidConfig(reasons []string) *ErrInvalidConfig {
	return &ErrInvalidConfig{Reasons: reasons}
}

func (e *ErrInvalidConfig) Error() string {
	return "invalid config: " + strings.Join(e.Reasons, ", ")
}

type ErrUnexpected struct {
	Err error
}

func NewErrUnexpected(err error) *ErrUnexpected {
	return &ErrUnexpected{Err: err}
}

func (e *ErrUnexpected) Error() string {
	return fmt.Sprintf("unexpected mongo error: %v", e.Err)
}

func (e *ErrUnexpected) Unwrap() error {
	return e.Err
}

type ErrNotFound struct {
}

func NewErrNotFound() *ErrNotFound {
	return &ErrNotFound{}
}

func (e *ErrNotFound) Error() string {
	return "document not found"
}

type ErrQueryInvalid struct {
	Reasons []string
}

func NewErrQueryInvalid(reasons []string) *ErrQueryInvalid {
	return &ErrQueryInvalid{Reasons: reasons}
}

func (e *ErrQueryInvalid) Error() string {
	return "invalid query: " + strings.Join(e.Reasons, ", ")
}

type ErrSortingInvalid struct {
	Reasons []string
}

func NewErrSortingInvalid(reasons []string) *ErrSortingInvalid {
	return &ErrSortingInvalid{Reasons: reasons}
}

func (e *ErrSortingInvalid) Error() string {
	return "sorting invalid: " + strings.Join(e.Reasons, ", ")
}
