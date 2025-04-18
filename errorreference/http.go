package errorreference

import (
	"errors"
	"net/http"
)

// errors related to http-based process activity
var (
	ErrorNotFound     = errors.New("not found")               // 404
	ErrorSlowDown     = errors.New("slow down")               // 429
	ErrorFailedToSend = errors.New("failed to send response") //500 case
	ErrInvalidRequest = errors.New("invalid request")
	ErrCuda700        = errors.New("got 700 from cuda invocation, will kill task")
	PanicDuringGoFunc = errors.New("paniced during a go func") // 400
)

var knownErrors = map[error]int{
	ErrorNotFound:     http.StatusNotFound,
	ErrorSlowDown:     http.StatusTooManyRequests,
	ErrorFailedToSend: http.StatusInternalServerError,
	ErrInvalidRequest: http.StatusBadRequest,
	//ErrCuda700: 500// TODO: ?
}

// TODO: special error types?
type HttpError interface {
	error
	statusCode() int // will return -1 for unknown
}

type httpErrorImpl struct {
	err  error
	code int
}

func (e httpErrorImpl) Error() string {
	return e.err.Error()
}

func (e httpErrorImpl) statusCode() int {
	return e.code
}

func WrapError(err error) httpErrorImpl { // TODO: ok? test?
	code := -1
	foundCode, exists := knownErrors[err]
	if exists {
		code = foundCode
	}
	return httpErrorImpl{
		err:  err,
		code: code,
	}
}
