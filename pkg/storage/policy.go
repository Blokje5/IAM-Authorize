package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

// Policy represents an authorization policy
// It determines whether User/Role/Service x is allowed to perform Action y on Resource z
type Policy struct {
	ID         int64 `json:"id"`
	audit      Audit
	Version    string      `json:"version"`
	Statements []Statement `json:"statements"`
}

const (
	Version = "v1"
)

// NewPolicy returns a new Policy
func NewPolicy(statements []Statement) *Policy {
	if statements == nil {
		statements = []Statement{}
	}

	return &Policy{
		Version:    Version,
		Statements: statements,
	}
}

// Scan implements sql.Scanner interface to parse policy directly
func (p *Policy) Scan(src interface{}) error {
	if src == nil {
		return ErrNotFound
	}
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed")
	}
	return json.Unmarshal(source, p)
}

// InsertPolicy inserts a new policy into the data store
func (s *Storage) InsertPolicy(ctx context.Context, policy *Policy) (*Policy, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin create policy tx: %w", err)
	}

	var ID int64
	err = tx.QueryRowContext(ctx, "SELECT insert_policy($1, $2, $3, $4, $5);", policy.Version, policy.audit.createdBy, policy.audit.lastModifiedBy, policy.audit.createdAt, policy.audit.lastModifiedAt).Scan(&ID)
	if err != nil {
		return nil, s.database.ProcessError(err)
	}
	policy.ID = ID

	for _, statement := range policy.Statements {
		_, err := tx.ExecContext(ctx, "SELECT insert_statement($1, $2, $3, $4);", ID, statement.Effect.String(), pq.Array(statement.Actions), pq.Array(statement.Resources))
		if err != nil {
			return nil, s.database.ProcessError(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("commit create policy tx: %w", err)
	}

	return policy, nil
}

// GetPolicy returns a policy based on the ID
func (s *Storage) GetPolicy(ctx context.Context, ID int64) (*Policy, error) {
	var policy Policy
	err := s.db.QueryRowContext(ctx, `SELECT query_policy($1);`, ID).Scan(&policy)
	if err != nil {
		fmt.Println(err)
		return nil, s.database.ProcessError(err)
	}

	return &policy, nil
}

func (s *Storage) ListPolicies(ctx context.Context) ([]Policy, error) {
	policies := []Policy{}
	rows, err := s.db.QueryContext(ctx, `SELECT query_policies();`)
	if err != nil {
		fmt.Println(err)
		return nil, s.database.ProcessError(err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			panic(err)
		}
	}()

	for rows.Next() {
		var policy Policy
		if err := rows.Scan(&policy); err != nil {
			return nil, err
		}
		policies = append(policies, policy)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return policies, nil
}

// DeletePolicy deletes a policy based on the ID
func (s *Storage) DeletePolicy(ctx context.Context, ID int64) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM policies WHERE id=$1;", ID)
	return s.database.ProcessError(err)
}

// Statement is a rule statement within a Policy
type Statement struct {
	Effect    Effect
	Actions   []string
	Resources []string
}

// NewStatement creates a new statement based on the given list of resources & actions
func NewStatement(effect Effect, resources, actions []string) Statement {
	if actions == nil {
		actions = []string{}
	}
	if resources == nil {
		resources = []string{}
	}

	return Statement{
		Effect:    effect,
		Actions:   actions,
		Resources: resources,
	}
}

// Effect = Allow/Deny
type Effect string

const (
	Deny  Effect = "deny"
	Allow Effect = "allow"
)

func (e Effect) String() string {
	return string(e)
}
