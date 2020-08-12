package jwt

type ErrJSONMarshalError struct {
	err error
}

func NewErrJSONMarshalError(err error) *ErrJSONMarshalError {
	return &ErrJSONMarshalError{err: err}
}

func (e *ErrJSONMarshalError) Error() string {
	return "jwt json marshal error: " + e.err.Error()
}

func (e *ErrJSONMarshalError) Unwrap() error {
	return e.err
}

type ErrSigningError struct {
	err error
}

func NewErrSigningError(err error) *ErrSigningError {
	return &ErrSigningError{err: err}
}

func (e *ErrSigningError) Error() string {
	return "jwt signing error: " + e.err.Error()
}

func (e *ErrSigningError) Unwrap() error {
	return e.err
}

type ErrJWTInvalid struct {
	err error
}

func NewErrJWTInvalid(err error) *ErrJWTInvalid {
	return &ErrJWTInvalid{err: err}
}

func (e *ErrJWTInvalid) Error() string {
	return "jwt invalid " + e.err.Error()
}

func (e *ErrJWTInvalid) Unwrap() error {
	return e.err
}

type ErrJWTVerificationFailure struct {
	err error
}

func NewErrJWTVerificationFailure(err error) *ErrJWTVerificationFailure {
	return &ErrJWTVerificationFailure{err: err}
}

func (e *ErrJWTVerificationFailure) Error() string {
	return "jwt verification failure " + e.err.Error()
}

func (e *ErrJWTVerificationFailure) Unwrap() error {
	return e.err
}
