package storage

import "context"

// User reprents a user with associated policies
// AuthN isn't handled by this server, so this should be fed
// by an external system
type User struct {
	ID    int64 `json:"id"`
	audit Audit
	Name  string `json:"name"`
}

// InsertUser inserts the User into the database and returns the User with id
func (s *Storage) InsertUser(ctx context.Context, user *User) (*User, error) {
	var ID int64
	err := s.db.QueryRowContext(ctx, "SELECT insert_user($1, $2, $3, $4, $5)", &user.Name, &user.audit.createdBy, &user.audit.createdBy, &user.audit.createdAt, &user.audit.createdAt).Scan(&ID)
	if err != nil {
		return nil, s.database.ProcessError(err)
	}

	user.ID = ID
	return user, nil
}

// ListUser returns all Users
func (s *Storage) ListUser(ctx context.Context) ([]User, error) {
	var users []User
	rows, err := s.db.QueryContext(ctx, `SELECT id, name, created_by, last_modified_by, created_at, last_modified_at FROM users;`)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			panic(err)
		}
	}()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.audit.createdBy, &user.audit.lastModifiedBy, &user.audit.createdAt, &user.audit.lastModifiedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Ensures we do not have to deal with nil pointers, instead we get an empty list
	if len(users) == 0 {
		users = []User{}
	}

	return users, nil
}

// GetUser returns a User based on the ID
func (s *Storage) GetUser(ctx context.Context, ID int64) (*User, error) {
	var user User
	err := s.db.QueryRowContext(ctx, `SELECT id, name, created_by, last_modified_by, created_at, last_modified_at FROM users WHERE id=$1;`, ID).Scan(&user.ID, &user.Name, &user.audit.createdBy, &user.audit.lastModifiedBy, &user.audit.createdAt, &user.audit.lastModifiedAt)
	if err != nil {
		return nil, s.database.ProcessError(err)
	}

	return &user, nil
}

// DeleteUser returns a User based on the ID
func (s *Storage) DeleteUser(ctx context.Context, ID int64) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM users WHERE id=$1;", ID)
	if err != nil {
		return s.database.ProcessError(err)
	}

	return nil
}
