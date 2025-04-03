package cuserrorx

import "errors"

const (
	Ok     = 0
	Failed = 2
)

type Error interface {
	error
	Code() int
	Msg() string
	GetError() error
}

type errorImpl struct {
	error
	code int
	msg  string
}

func NewWithError(code int, err error) Error {
	return &errorImpl{
		error: err,
		code:  code,
	}
}

func NewWithErrorMsg(msg string) Error {
	return &errorImpl{
		code: Failed,
		msg:  msg,
	}
}

func NewWithMsg(code int, msg string) Error {
	return &errorImpl{
		code: code,
		msg:  msg,
	}
}

func NewWithCode(code int) Error {
	return &errorImpl{
		code: code,
	}
}

func (e *errorImpl) Error() string {
	if e.error != nil {
		return e.error.Error()
	}

	return e.Msg()
}

func (e *errorImpl) Code() int {
	return e.code
}

func (e *errorImpl) Msg() string {
	if e.msg != "" {
		return e.msg
	}

	return ""
}

func (e *errorImpl) GetError() error {
	return e.error
}

func As(err error) Error {
	se := &errorImpl{}
	if !errors.As(err, &se) {
		return nil
	}

	return se
}

func Is(err, target error) bool {
	if errors.Is(err, target) {
		return true
	}

	e := As(err)
	if e == nil || e.GetError() == nil {
		return false
	}

	return Is(e.GetError(), target)
}
