package helpers

import (
	"net/http"
	"robinhood-assignment/internal/dto"
)

func ErrorHandler(err error) *dto.ErrorResponse {
	errRes := dto.ErrorResponse{}
	if e, ok := err.(CustomError); ok {
		errRes.StatusCode = e.Code()
		errRes.Error = e.Error()
		return &errRes
	}
	errRes.StatusCode = http.StatusInternalServerError
	errRes.Error = err.Error()
	return &errRes
}
