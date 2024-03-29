package storage

import "context"

// Namespace represents a scoped authorization namespace
// E.g. the Namespace users would be the namespace for authorization requests related to users
type Namespace struct {
	ID    int64 `json:"id"`
	audit Audit
	Name  string `json:"name"`
}

// NamespaceList represents a list of namespaces
type NamespaceList struct {
	Namespaces []Namespace `json:"namespaces"`
}

// ListNamespaces returns a list of all namespaces available to the user
func (s *Storage) ListNamespaces(ctx context.Context) (*NamespaceList, error) {
	var namespaces []Namespace
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
		namespaces = append(namespaces, namespace)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Ensures we do not have to deal with nil pointers, instead we get an empty list
	if len(namespaces) == 0 {
		namespaces = []Namespace{}
	}

	return &NamespaceList{Namespaces: namespaces}, nil
}

// GetNamespace returns a namespace based on the ID
func (s *Storage) GetNamespace(ctx context.Context, ID int64) (*Namespace, error) {
	var namespace Namespace
	err := s.db.QueryRowContext(ctx, "SELECT id, name, created_by, last_modified_by, created_at, last_modified_at FROM namespaces WHERE id=$1;", ID).Scan(&namespace.ID, &namespace.Name, &namespace.audit.createdBy, &namespace.audit.lastModifiedBy, &namespace.audit.createdAt, &namespace.audit.lastModifiedAt)
	if err != nil {
		return nil, s.database.ProcessError(err)
	}

	return &namespace, nil
}

// InsertNamespace inserts the namespace into the database and returns the namespace with id
func (s *Storage) InsertNamespace(ctx context.Context, namespace Namespace) (*Namespace, error) {
	var ID int64
	err := s.db.QueryRowContext(ctx, "INSERT INTO namespaces (name, created_by, last_modified_by, created_at, last_modified_at) VALUES ($1, $2, $3, $4, $5) RETURNING id;", namespace.Name, namespace.audit.createdBy, namespace.audit.createdBy, namespace.audit.createdAt, namespace.audit.createdAt).Scan(&ID)
	if err != nil {
		return nil, s.database.ProcessError(err)
	}

	namespace.ID = ID
	return &namespace, nil
}

// UpdateNamespace updates the namespace (if there are changes) and returns the updated namespace object
func (s *Storage) UpdateNamespace(ctx context.Context, namespace Namespace) (*Namespace, error) {
	var ID int64
	err := s.db.QueryRowContext(ctx, "UPDATE namespaces SET name = $1, last_modified_by = $2, last_modified_at = $3 WHERE id = $4 RETURNING id;", namespace.Name, namespace.audit.lastModifiedBy, namespace.audit.lastModifiedAt, namespace.ID).Scan(&ID)
	if err != nil {
		return nil, s.database.ProcessError(err)
	}

	namespace.ID = ID
	return &namespace, nil
}

// DeleteNamespace deletes the namespace based on the ID
func (s *Storage) DeleteNamespace(ctx context.Context, ID int64) {
	s.db.ExecContext(ctx, "DELETE FROM namespace WHERE id=$1;", ID)
}
