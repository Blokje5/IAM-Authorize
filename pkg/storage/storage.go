package storage

import "database/sql"

import "github.com/blokje5/iam-server/pkg/storage/database"

// Storage wraps the backend database
// and implements all  methods necessary for storage
type Storage struct {
	db *sql.DB
	database database.Database
}

// New returns a new instance of the Storage with the given db as backend
func New(db *sql.DB, database database.Database) *Storage {
	return &Storage{
		db: db,
		database: database,
	}
}
