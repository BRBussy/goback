package authentication

type ErrLoginFailed struct {
}

func NewErrLoginFailed() *ErrLoginFailed {
	return &ErrLoginFailed{}
}

func (e *ErrLoginFailed) Error() string {
	return "login failed: incorrect email address or password"
}
