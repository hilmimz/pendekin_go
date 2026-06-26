package errs

import "net/http"

type Error struct {
	Code    int
	Message string
	Err     error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func NotFound(msg string, err error) *Error {
	return &Error{
		Code:    http.StatusNotFound,
		Message: msg,
		Err:     err,
	}
}

func Conflict(msg string, err error) *Error {
	return &Error{
		Code:    http.StatusConflict,
		Message: msg,
		Err:     err,
	}
}

func Internal(msg string, err error) *Error {
	return &Error{
		Code:    http.StatusInternalServerError,
		Message: msg,
		Err:     err,
	}
}

func Unauthorized(msg string, err error) *Error {
	return &Error{
		Code:    http.StatusUnauthorized,
		Message: msg,
		Err:     err,
	}
}

func Forbidden(msg string, err error) *Error {
	return &Error{
		Code:    http.StatusForbidden,
		Message: msg,
		Err:     err,
	}
}
