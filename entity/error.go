package entity

import (
	"fmt"

	"github.com/pkg/errors"
)

var ErrNotValidCategory = errors.New("Category is not valid")

var ErrCategoryHasWrongLevel = errors.New("Level is wrong")

//ErrNotFound not found
var ErrNotFound = errors.New("Not found")

//ErrInvalidEntity invalid entity
var ErrInvalidEntity = errors.New("Invalid entity")

//ErrCannotBeDeleted cannot be deleted
var ErrCannotBeDeleted = errors.New("Cannot Be Deleted")

//ErrNotEnoughBooks cannot borrow
var ErrNotEnoughBooks = errors.New("Not enough books")

//ErrBookAlreadyBorrowed cannot borrow
var ErrBookAlreadyBorrowed = errors.New("Book already borrowed")

//ErrBookNotBorrowed cannot return
var ErrBookNotBorrowed = errors.New("Book not borrowed")

type ErrorType struct {
	t string
}

var (
	ErrorTypeUnknown        = ErrorType{"unknown"}
	ErrorTypeAuthorization  = ErrorType{"authorization"}
	ErrorTypeIncorrectInput = ErrorType{"incorrect-input"}
	ErrorTypeDatabase       = ErrorType{"database"}
)

type QueryError struct {
	Query string
	Err   error
}

func (e *QueryError) Error() string {
	return e.Query + ": " + e.Err.Error()
}

type AppError struct {
	Err  error
	Msg  string
	Type ErrorType
}

func (err AppError) Error() string {

	type causer interface {
		Cause() error
	}

	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	// fmt.Printf("\n\n cause : \n %+v \n\n", err.Err.(causer).Cause().Error())
	// fmt.Printf("\n\n st : \n %+v \n\n", errors.Cause(err.Err).(stackTracer).StackTrace())

	//	return fmt.Sprintf("error caused due to %v; message: %v; additional context: %v", e.Inner, e.Message, e.AdditionalContext)

	return fmt.Sprintf("%+v", err.Err)
}

// // Unwrap is used to make it work with errors.Is, errors.As.
// func (e *ErrZeroDivision) Unwrap() error {
// 	// Return the inner error.
// 	return e.Inner
// }

func NewAppError(err error) AppError {
	return AppError{
		Err:  err,
		Type: ErrorTypeIncorrectInput,
	}
}

func NewAppDbError(err error) AppError {
	return AppError{
		Err:  err,
		Type: ErrorTypeDatabase,
	}
}

func CustomeErrorInstance() error {
	return errors.New("File type not supported")
}
