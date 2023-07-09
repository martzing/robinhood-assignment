package helpers

import "net/http"

type CustomError struct {
	StatusCode int
	Message    string
}

func NewCustomError(code int, message string) error {
	return CustomError{
		StatusCode: code,
		Message:    message,
	}
}

func (e CustomError) Error() string {
	return e.Message
}

func (e CustomError) Code() int {
	return e.StatusCode
}

var InternalError = NewCustomError(http.StatusInternalServerError, "Something went wrong please contact developer.")
