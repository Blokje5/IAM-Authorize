package storage

import (
	"context"

)

func (s *Storage) InsertPolicyForUser(ctx context.Context, userID int64, policyID int64) (error) {
	_, err := s.db.ExecContext(ctx, `INSERT INTO users_policies (user_id, policy_id) VALUES ($1, $2);`, userID, policyID)
	if err != nil {
		return s.database.ProcessError(err)
	}

	return nil
}