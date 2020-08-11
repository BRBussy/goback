package key

type ErrNilKey struct {
}

func NewErrNilKey() *ErrNilKey {
	return &ErrNilKey{}
}

func (e *ErrNilKey) Error() string {
	return "key nil"
}

type ErrParsingError struct {
	err error
}

func NewErrParsingError(err error) *ErrParsingError {
	return &ErrParsingError{err: err}
}

func (e *ErrParsingError) Error() string {
	return "parsing error: " + e.err.Error()
}

func (e *ErrParsingError) Unwrap() error {
	return e.err
}

type ErrGenerationError struct {
	err error
}

func NewErrGenerationError(err error) *ErrGenerationError {
	return &ErrGenerationError{err: err}
}

func (e *ErrGenerationError) Error() string {
	return "generation error: " + e.err.Error()
}

func (e *ErrGenerationError) Unwrap() error {
	return e.err
}
