package authentication

type ErrLoginFailed struct {
}

func NewErrLoginFailed() *ErrLoginFailed {
	return &ErrLoginFailed{}
}

func (e *ErrLoginFailed) Error() string {
	return "login failed: incorrect email address or password"
}

type ErrJWTInvalid struct {
}

func NewErrJWTInvalid() *ErrJWTInvalid {
	return &ErrJWTInvalid{}
}

func (e *ErrJWTInvalid) Error() string {
	return "jwt invalid"
}

type ErrJWTExpired struct{}

func NewErrJWTExpired() *ErrJWTExpired {
	return &ErrJWTExpired{}
}

func (e *ErrJWTExpired) Error() string {
	return "jwt expired"
}

type ErrJSONUnmarshalError struct {
	err error
}

func NewErrJSONUnmarshalError(err error) *ErrJSONUnmarshalError {
	return &ErrJSONUnmarshalError{err: err}
}

func (e *ErrJSONUnmarshalError) Error() string {
	return "jwt json payload unmarshal error: " + e.err.Error()
}

func (e *ErrJSONUnmarshalError) Unwrap() error {
	return e.err
}
