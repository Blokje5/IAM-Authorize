package storage

import "context"

// Namespace represents a scoped authorization namespace
// E.g. the Namespace users would be the namespace for authorization requests related to users
type Namespace struct {
	ID    int `json:"id"`
	audit Audit
	Name  string `json:"name"`
}

// ListNamespaces returns a list of all namespaces available to the user
func (s *Storage) ListNamespaces(ctx context.Context) ([]*Namespace, error) {
	var namespaces []*Namespace
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, name, created_by, last_modified_by, created_at, last_modified_at from namespaces")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			panic(err)
		}
	}()

	for rows.Next() {
		var namespace Namespace
		if err := rows.Scan(&namespace.ID, &namespace.Name, &namespace.audit.createdBy, &namespace.audit.lastModifiedBy, &namespace.audit.createdAt, &namespace.audit.lastModifiedAt); err != nil {
			return nil, err
		}
		namespaces = append(namespaces, &namespace)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return namespaces, nil
}

// GetNamespace returns a namespace based on the ID
func (s *Storage) GetNamespace(ctx context.Context, ID int) *Namespace {
	var namespace Namespace
	s.db.QueryRowContext(ctx, "SELECT id, name, created_by, last_modified_by, created_at, last_modified_at FROM namespaces WHERE id=$1;", ID).Scan(&namespace.ID, &namespace.Name, &namespace.audit.createdBy, &namespace.audit.lastModifiedBy, &namespace.audit.createdAt, &namespace.audit.lastModifiedAt)

	return &namespace
}

// InsertNamespace inserts the namespace into the database and returns the namespace with id
func (s *Storage) InsertNamespace(ctx context.Context, namespace Namespace) *Namespace {
	var ID int
	s.db.QueryRowContext(ctx, "INSERT INTO namespaces (name, created_by, last_modified_by, created_at, last_modified_at) VALUES ($1, $2, $3, $4, $5) RETURNING id;", namespace.Name, namespace.audit.createdBy, namespace.audit.createdBy, namespace.audit.createdAt, namespace.audit.createdAt).Scan(&ID)

	namespace.ID = ID
	return &namespace
}

// UpdateNamespace updates the namespace (if there are changes) and returns the updated namespace object
func (s *Storage) UpdateNamespace(ctx context.Context, namespace Namespace) *Namespace {
	var ID int
	s.db.QueryRowContext(ctx, "UPDATE namespaces SET (name, last_modified_by, last_modified_at) VALUES ($1, $2, $3) WHERE id = $4 RETURNING id;", namespace.Name, namespace.audit.lastModifiedBy, namespace.audit.lastModifiedAt, namespace.ID).Scan(&ID)

	namespace.ID = ID
	return &namespace
}

// DeleteNamespace deletes the namespace based on the ID
func (s *Storage) DeleteNamespace(ctx context.Context, ID int) {
	s.db.ExecContext(ctx, "DELETE FROM namespace WHERE id=$1;", ID)
}
