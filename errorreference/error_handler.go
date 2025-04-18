package errorreference

type ErrorHandler struct {
	Message    string
	StatusCode int
	Err        error
}

func (e *ErrorHandler) SetError(message string, statusCode int) {
	e.Message = message
	e.StatusCode = statusCode
}
