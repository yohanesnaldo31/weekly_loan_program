package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	userLoanKey = "loan_user_%d"
	userLoanTTL = 10 * time.Minute
)

// GetUserLoanByUserID returns the cached loans for the given user, stored
// under key "loan_user_<user_id>". It returns a redis.Nil-wrapped error when
// the key does not exist.
func (r *Repository) GetUserLoanByUserID(ctx context.Context, userID int64) ([]Loan, error) {
	key := fmt.Sprintf(userLoanKey, userID)
	value, err := r.redis.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("cache: get loan for user_id %d: %w", userID, err)
	}

	var loans []Loan
	if err := json.Unmarshal(value, &loans); err != nil {
		return nil, fmt.Errorf("cache: unmarshal loan for user_id %d: %w", userID, err)
	}

	return loans, nil
}

// SetUserLoanByUserID caches loans under key "loan_user_<user_id>" as a JSON
// string, expiring after 10 minutes.
func (r *Repository) SetUserLoanByUserID(ctx context.Context, userID int64, loans []Loan) error {
	value, err := json.Marshal(loans)
	if err != nil {
		return fmt.Errorf("cache: marshal loan for user_id %d: %w", userID, err)
	}

	key := fmt.Sprintf(userLoanKey, userID)
	if err := r.redis.Set(ctx, key, value, userLoanTTL).Err(); err != nil {
		return fmt.Errorf("cache: set loan for user_id %d: %w", userID, err)
	}

	return nil
}

// DeleteUserLoanByUserID removes the cached loans for the given user, stored
// under key "loan_user_<user_id>".
func (r *Repository) DeleteUserLoanByUserID(ctx context.Context, userID int64) error {
	key := fmt.Sprintf(userLoanKey, userID)
	if err := r.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("cache: delete loan for user_id %d: %w", userID, err)
	}

	return nil
}
