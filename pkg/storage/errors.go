package storage

import "errors"

var (
	// ErrUniqueViolation is the error returned when a Uniqueness Constraint is violated in the database
	ErrUniqueViolation = errors.New("Unique constraint violation")
	// ErrDatabase handles uncaught database exceptions
	ErrDatabase = errors.New("Internal Database Error")
)
