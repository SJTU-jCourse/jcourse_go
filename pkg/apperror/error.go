package apperror

import "net/http"

type AppError struct {
	Code       int
	Msg        string
	Err        error
	HTTPStatus int
}

func (e *AppError) Error() string {
	return e.Msg
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) StatusCode() int {
	return e.HTTPStatus
}

func (e *AppError) Wrap(err error) *AppError {
	return &AppError{
		Code:       e.Code,
		Msg:        e.Msg,
		Err:        err,
		HTTPStatus: e.HTTPStatus,
	}
}

var (
	ErrSomethingWrong = &AppError{Code: -1000, Msg: "Something went wrong", HTTPStatus: http.StatusInternalServerError}
	ErrNotFound       = &AppError{Code: -1001, Msg: "Not Found", HTTPStatus: http.StatusNotFound}
)
