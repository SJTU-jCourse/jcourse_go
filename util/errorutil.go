package util

type NotImplementError struct {
	Msg string
}

func (e NotImplementError) Error() string {
	return e.Msg
}
func NewNotImplementError(msg string) *NotImplementError {
	return &NotImplementError{Msg: msg}
}

var ErrNotImplemented error = &NotImplementError{Msg: "Not implemented yet"}

type InvalidParamError struct {
	Msg string
}

func (e InvalidParamError) Error() string {
	return e.Msg
}
func NewInvalidParamError(msg string) *InvalidParamError {
	return &InvalidParamError{Msg: msg}
}

var _ error = &InvalidParamError{Msg: "参数错误"}

type InternalServerError struct {
	Msg string
}

func (e InternalServerError) Error() string {
	return e.Msg
}
func NewInternalServerError(msg string) *InternalServerError {
	return &InternalServerError{Msg: msg}
}

var ErrInternal error = &InternalServerError{"内部错误"}

type NotFoundError struct {
	Msg string
}

func (e NotFoundError) Error() string {
	return e.Msg
}
func NewNotFoundError(msg string) *NotFoundError {
	return &NotFoundError{Msg: msg}
}

var ErrNotFound error = &NotFoundError{Msg: "资源不存在"}
