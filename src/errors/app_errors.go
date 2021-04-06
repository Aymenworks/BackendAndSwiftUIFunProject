package errors

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
)

type appError struct {
	Code           string
	Message        string
	HTTPStatusCode int
}

func (e *appError) Error() string {
	return fmt.Sprintf("%v %v", e.Code, e.Message)
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
	if err == nil {
		return newUnknownError()
	}

	var appErr *appError

	if errors.As(err, &appErr) {
		return appErr
	}

	return newUnknownError()
}

func newBadRequest(code, msg string) *appError {
	return &appError{
		Message:        msg,
		HTTPStatusCode: http.StatusBadRequest,
	}
}

func newInternalError(code, msg string) *appError {
	return &appError{
		Message:        msg,
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func newNotFoundError() *appError {
	return &appError{
		HTTPStatusCode: http.StatusNotFound,
	}
}

func newUnknownError() *appError {
	return &appError{
		Message:        "Unknown error",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func Wrap(err error, msg string) error {
	if err == nil {
		return err
	}

	pc := make([]uintptr, 1)
	// Skip the first 2 because the first one is the internal system and the second is this method itself, but we are interested by the ones just before
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	s := fmt.Sprintf("%v %v:%d %v\n", msg, frame.File, frame.Line, frame.Function)

	return fmt.Errorf("%w %s", err, s)
}
