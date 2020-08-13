package claims

type ErrJSONUnmarshallError struct {
	Err error
}

func NewErrJSONUnmarshallError(err error) *ErrJSONUnmarshallError {
	return &ErrJSONUnmarshallError{Err: err}
}

func (e *ErrJSONUnmarshallError) Error() string {
	return "claims json unmarshall error: " + e.Err.Error()
}

func (e *ErrJSONUnmarshallError) Unwrap() error {
	return e.Err
}

type ErrClaimsNotInContext struct{}

func NewErrClaimsNotInContext() *ErrClaimsNotInContext {
	return &ErrClaimsNotInContext{}
}

func (e *ErrClaimsNotInContext) Error() string {
	return "claims not in context"
}
