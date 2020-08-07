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

type ErrInvalidType struct {
	Type Type
}

func NewErrInvalidType(t Type) *ErrInvalidType {
	return &ErrInvalidType{Type: t}
}

func (e *ErrInvalidType) Error() string {
	return "invalid claims type: " + e.Type.String()
}
