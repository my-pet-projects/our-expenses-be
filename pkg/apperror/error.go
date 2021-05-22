package apperror

import (
	"fmt"

	"github.com/pkg/errors"
)

// ErrorType holds error type.
type ErrorType struct {
	t string
}

var (
	// ErrorTypeUnknown is unknown error.
	ErrorTypeUnknown = ErrorType{"unknown"}
	// ErrorTypeAuthorization is authorization error.
	ErrorTypeAuthorization = ErrorType{"authorization"}
	// ErrorTypeIncorrectInput is incorrect input error.
	ErrorTypeIncorrectInput = ErrorType{"incorrect-input"}
	// ErrorTypeDatabase is database error.
	ErrorTypeDatabase = ErrorType{"database"}
)

// AppError is custom error with application specific fields.
type AppError struct {
	Err  error
	Msg  string
	Type ErrorType
}

func (err AppError) Error() string {
	return fmt.Sprintf("%+v", err.Err)
}

// GetStackTrace returns error stack trace if error implements stackTracer interface.
func GetStackTrace(err error) string {
	if err == nil {
		return ""
	}

	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	e, ok := err.(stackTracer)
	if ok {
		st := e.StackTrace()
		return fmt.Sprintf("%+v", st)
	}

	return ""
}

// GetCause returns underlying error cause if error implements causer interface, otherwise original error.
func GetCause(err error) string {
	if err == nil {
		return ""
	}

	type causer interface {
		Cause() error
	}

	e, ok := err.(causer)
	if ok {
		return e.Cause().Error()
	}

	return err.Error()
}
