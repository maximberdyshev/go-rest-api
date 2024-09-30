package errs

import "encoding/json"

var (
	ErrBadRequest    = NewAppError(nil, "bad request")
	ErrUnauthorized  = NewAppError(nil, "unauthorized")
	ErrNotFound      = NewAppError(nil, "not found")
	ErrInternal      = NewAppError(nil, "internal server error")
	ErrIncorrectBody = NewAppError(nil, "incorrect body")
)

type AppError struct {
	Err error
	Msg string
}

func NewAppError(err error, msg string) *AppError {
	return &AppError{
		Err: err,
		Msg: msg,
	}
}

func (e *AppError) Error() string {
	return e.Msg
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Marshal() []byte {
	if m, err := json.Marshal(e); err == nil {
		return m
	}
	return nil
}
