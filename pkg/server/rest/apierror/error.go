package apierror

import (
	"fmt"
)

type Error struct {
	StatusCode int   `json:"statusCode"`
	Msg        any   `json:"msg"`
	err        error // the original error. if the error cause exists - it will be logged
}

func New(statusCode int, msg any) error {
	return &Error{
		StatusCode: statusCode,
		Msg:        msg,
	}
}

func NewWithError(statusCode int, msg string, err error) error {
	if err == nil {
		return nil
	}

	return Error{
		StatusCode: statusCode,
		Msg:        msg,
		err:        err,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("%v: %v", e.Msg, e.err)
}

func (err *Error) Unwrap() error { return err.err }
