package errors

import (
	"fmt"
	"net/http"

	"golang.org/x/xerrors"
)

type appError struct {
	Code           string
	Message        string
	HTTPStatusCode int
}

type AppError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%v %v", e.Code, e.Message)
}

func (e *appError) New() *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message,
	}
}

func AsAppError(err error) *appError {
	var appErr *appError
	if xerrors.As(err, appErr) {
		return appErr
	}

	return nil
}

func NewBadRequest(code, msg string) *appError {
	return &appError{
		Message:        msg,
		HTTPStatusCode: http.StatusBadRequest,
	}
}

func NewInternalError(code, msg string) *appError {
	return &appError{
		Message:        msg,
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func NewNotFoundError() *appError {
	return &appError{
		HTTPStatusCode: http.StatusNotFound,
	}
}

func NewUnknownError() *appError {
	return &appError{
		Message:        "Unknown error",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}
