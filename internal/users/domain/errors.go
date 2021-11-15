package domain

import "github.com/pkg/errors"

// Errors.
var (
	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("password or user is incorrect")
)
