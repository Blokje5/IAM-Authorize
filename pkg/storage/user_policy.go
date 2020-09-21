package storage

import (
	"context"
)

func (s *Storage) InsertPolicyForUser(ctx context.Context, userID int64, policyID int64) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO users_policies (user_id, policy_id) VALUES ($1, $2);`, userID, policyID)
	if err != nil {
		return s.database.ProcessError(err)
	}

	return nil
}

func (s *Storage) GetPolicyIDsForUser(ctx context.Context, user *User) ([]int64, error) {
	var policies []int64
	rows, err := s.db.QueryContext(ctx, `SELECT query_policies_for_user($1);`, &user.Name)
	if err != nil {
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
		policies = append(policies, policy.ID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return policies, nil
}

func (s *Storage) GetPoliciesForUser(ctx context.Context, user *User) ([]Policy, error) {
	policies := []Policy{}
	rows, err := s.db.QueryContext(ctx, `SELECT query_policies_for_user($1);`, &user.Name)
	if err != nil {
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
