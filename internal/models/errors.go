package models

import "errors"

var (
	ErrDatabase        = errors.New("database error")
	ErrNotFound        = errors.New("not found")
	ErrInternal        = errors.New("internal error")
	ErrInvalid         = errors.New("invalid input")
	ErrCanceled        = errors.New("canceled")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrTooManyRequests = errors.New("too many requests")
	ErrConflicted      = errors.New("conflicted")
	ErrMalformed       = errors.New("malformed")
)
