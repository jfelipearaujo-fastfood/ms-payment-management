package custom_error

type BusinessError struct {
	code    int
	title   string
	message string
}

func New(code int, title string, message string) BusinessError {
	return BusinessError{
		code:    code,
		title:   title,
		message: message,
	}
}

func (e BusinessError) Code() int {
	return e.code
}

func (e BusinessError) Title() string {
	return e.title
}

func (e BusinessError) Error() string {
	return e.message
}

func IsBusinessErr(err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(BusinessError)
	return ok
}
