package database

import (
	"context"
	"database/sql"
)

// Database is an interface that should be implemented by the underlying database struct.
// It provides methods for initialization and error handling
type Database interface {
	ProcessError(error) error
	Initialize(context.Context) (*sql.DB, error)
}
