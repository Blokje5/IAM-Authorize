package storage

import "errors"

var (
	// ErrUniqueViolation is the error returned when a Uniqueness Constraint is violated in the database
	ErrUniqueViolation = errors.New("Unique constraint violation")
	// ErrForeignKeyViolation is returned when a foreign key is violated in the database
	ErrForeignKeyViolation = errors.New("Foreign key violation")
	// ErrDatabase handles uncaught database exceptions
	ErrDatabase = errors.New("Internal Database Error")
	// ErrNotFound handles empty data
	ErrNotFound = errors.New("Did not found resource")
)
