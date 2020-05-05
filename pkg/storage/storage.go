package storage

import "database/sql"

// Storage wraps the backend database
// and implements all  methods necessary for storage
type Storage struct {
	db *sql.DB
}

// New returns a new instance of the Storage with the given db as backend
func New(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}