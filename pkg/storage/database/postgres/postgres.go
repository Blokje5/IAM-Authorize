package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"net/url"
	"time"

	"github.com/blokje5/iam-server/pkg/storage"
	"github.com/cenkalti/backoff/v4"

	"github.com/golang-migrate/migrate/v4"
	// Migration Driver
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	// Used as SQL driver
	_ "github.com/jackc/pgx"
)

// Postgres implements the DB interface and wraps a PG sql.DB instance
type Postgres struct {
	config *PostgresConfig
}

// New returns a new reference to a Postgres Instance
func New(config *PostgresConfig) *Postgres {
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

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 10 * time.Second
	op := func() error {
		return db.PingContext(ctx)
	}
	err = backoff.Retry(op, b)
	if err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	m, err := migrate.New(
		addFileScheme(p.config.MigrationPath),
		p.config.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("error fetching migration configuration: %w", err)
	}

	_, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return nil, fmt.Errorf("error fetching migration version: %w", err)
	}

	if dirty {
		return nil, fmt.Errorf("database in invalid state after previous migration")
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("database migration failed: %w", err)
	}

	return db, nil
}

// ProcessError handles any error thrown in the sql layer and checks if it is a known postgres error
func (p *Postgres) ProcessError(err error) error {
	if err, ok := err.(*pq.Error); ok {
		return p.handlePGError(err)
	}

	return err
}

func (p *Postgres) handlePGError(e *pq.Error) error {
	switch e.Code {
	case "23505": // Unique violation
		return storage.ErrUniqueViolation
	}

	return storage.ErrDatabase
}

func addFileScheme(p string) string {
	u := url.URL{}
	u.Scheme = "file"
	u.Path = p
	return u.String()
}
