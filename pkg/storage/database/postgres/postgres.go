package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// Used as SQL driver
	_ "github.com/jackc/pgx"
)

// Postgres implements the DB interface and wraps a PG sql.DB instance
type Postgres struct {
	config PostgresConfig
}

// New returns a new reference to a Postgres Instance
func New(config PostgresConfig) *Postgres {
	p := &Postgres{
		config: config,
	}

	return p
}

// Initialize initializes the postgres instance, validating the connection
// and performing any migration logic necessary
func (p *Postgres) Initialize(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("postgres", p.config.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("error returning database driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres", driver)

	version, dirty, err := m.Version()
	if err != nil {
		return nil, fmt.Errorf("error fetching migration version: %w", err)
	}

	if dirty {
		return nil, fmt.Errorf("database in invalid state after previous migration")
	}

	err = m.Migrate(version)
	if err != nil {
		return nil, fmt.Errorf("database migration failed: %w", err)
	}

	return db, nil
}
