package filter

type ErrInvalidType struct {
	Type Type
}

func NewErrInvalidType(t Type) *ErrInvalidType {
	return &ErrInvalidType{Type: t}
}

func (e *ErrInvalidType) Error() string {
	return "invalid filter type: " + e.Type.String()
}

type ErrJSONUnmarshallError struct {
	Err error
}

func NewErrJSONUnmarshallError(err error) *ErrJSONUnmarshallError {
	return &ErrJSONUnmarshallError{Err: err}
}

func (e *ErrJSONUnmarshallError) Error() string {
	return "filter json unmarshall error: " + e.Err.Error()
}

func (e *ErrJSONUnmarshallError) Unwrap() error {
	return e.Err
}
