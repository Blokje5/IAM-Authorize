package storage

import (
	"context"
	"database/sql"
	"github.com/blokje5/iam-server/pkg/storage/database"
)

// Storage wraps the backend database
// and implements all  methods necessary for storage
type Storage struct {
	db       *sql.DB
	database database.Database
}

// New returns a new instance of the Storage with the given db as backend
func New(db *sql.DB, database database.Database) *Storage {
	return &Storage{
		db:       db,
		database: database,
	}
}

// Clean cleans everything in the storage. This is a
// destructive operation and should only be used internally
// It does not guarantee the database is in a clean state
func (s *Storage) Clean(ctx context.Context) error {
	// TODO wrap in transaction
	_, err := s.db.ExecContext(ctx, "ALTER SEQUENCE namespaces_id_seq RESTART WITH 1;")
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, "TRUNCATE TABLE namespaces;")
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, "TRUNCATE TABLE users_policies;")
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, "ALTER SEQUENCE users_id_seq RESTART WITH 1;")
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, "TRUNCATE TABLE users CASCADE;")
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, "TRUNCATE TABLE statements;")
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, "ALTER SEQUENCE policies_id_seq RESTART WITH 1;")
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, "TRUNCATE TABLE policies CASCADE;")
	if err != nil {
		return err
	}

	return err
}
