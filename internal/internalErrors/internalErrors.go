package internalErrors

type ErrorIncorrectNumberArgs struct{}

func NewErrorIncorrectNumberArgs() *ErrorIncorrectNumberArgs{
	return &ErrorIncorrectNumberArgs{}
}

func (e *ErrorIncorrectNumberArgs)Error() string{
	return "incorrect number args"
}

type ErrorAuthentication struct {
	errAuth error
}

func NewErrorAuthentication(err error) *ErrorAuthentication{
	return &ErrorAuthentication{errAuth: err}
}

func (e *ErrorAuthentication)Error() string{
	return "Error authentication" + e.errAuth.Error()
}

type ErrorIncorrectPassword struct {
}

func NewErrorIncorrectPassword() *ErrorIncorrectPassword{
	return &ErrorIncorrectPassword{}
}

func (e *ErrorIncorrectPassword)Error() string{
	return "incorrect password"
}