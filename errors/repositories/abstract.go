package errorsRepositories

const GeneralError = 1
const ApiError = 2
const AuthError = 3
const NotFoundError = 4

type Interface interface {
	SetError(text string)
	SetCode(code int)
	Code() int

	error
}

type Error struct {
	text string
	code int

	Interface
}

func (e *Error) SetError(text string) {
	e.text = text
}

func (e *Error) Error() string {
	return e.text
}

func (e *Error) SetCode(code int) {
	e.code = code
}

func (e *Error) Code() int {
	return e.code
}

func NewError(text string, code int) *Error {
	return &Error{
		text: text,
		code: code,
	}
}
