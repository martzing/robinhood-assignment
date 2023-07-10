package helpers

import (
	"net/http"
	"robinhood-assignment/internal/dto"
)

type customError struct {
	StatusCode int
	Message    string
}

func NewCustomError(code int, message string) error {
	return customError{
		StatusCode: code,
		Message:    message,
	}
}

func (e customError) Error() string {
	return e.Message
}

func (e customError) Code() int {
	return e.StatusCode
}

func ErrorHandler(err error) *dto.ErrorResponse {
	errRes := dto.ErrorResponse{}
	if e, ok := err.(customError); ok {
		errRes.StatusCode = e.Code()
		errRes.Error = e.Error()
		return &errRes
	}
	errRes.StatusCode = http.StatusInternalServerError
	errRes.Error = err.Error()
	return &errRes
}

var InternalError = NewCustomError(http.StatusInternalServerError, "Something went wrong please contact developer.")
