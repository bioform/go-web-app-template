package errors

import "errors"

var (
	// ErrRecordNotFound is used when a user is not found in the database
	ErrRecordNotFound = errors.New("record not found")
)
