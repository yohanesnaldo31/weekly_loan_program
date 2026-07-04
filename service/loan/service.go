package loan

import (
	"context"
	"time"

	"weekly_loan_program/repo/cache"
	"weekly_loan_program/repo/db"

	"github.com/jackc/pgx/v5"
)

type dbRepoProvider interface {
	// Begin starts a transaction. The caller must call [pgx.Tx.Commit] or
	// [pgx.Tx.Rollback] on the returned transaction to finalize it.
	Begin(ctx context.Context) (pgx.Tx, error)

	// GetUserLoansByUserID returns up to 10 of the given user's loans, sorted by
	// create_time descending (most recent first).
	GetUserLoansByUserID(ctx context.Context, userID int64) ([]db.Loan, error)

	// GetBillingsByLoanIDAndDueDate returns the billings for the given loan due
	// before dueDate, sorted by due_date ascending. When status is 0, billings
	// are returned regardless of status; otherwise only billings matching status
	// are returned.
	GetBillingsByLoanIDAndDueDate(ctx context.Context, loanID int64, dueDate time.Time, status int16) ([]db.Billing, error)

	// InsertBilling inserts billing within the given transaction and returns the
	// generated billing ID.
	InsertBilling(ctx context.Context, tx pgx.Tx, billing db.Billing) (int64, error)

	// InsertLoan inserts loan within the given transaction and returns the
	// generated loan ID.
	InsertLoan(ctx context.Context, tx pgx.Tx, loan db.Loan) (int64, error)

	// UpdateBillings updates the status and payment time of the billings with
	// the given IDs within the given transaction.
	UpdateBillings(ctx context.Context, tx pgx.Tx, billingIDs []int64, status int16, paymentTime *time.Time) error

	// UpdateLoan updates the total paid, status, and update time of the loan
	// with the given ID within the given transaction.
	UpdateLoan(ctx context.Context, tx pgx.Tx, loanID int64, totalPaid int64, status int16, updateTime time.Time) error
}

type cacheRepoProvider interface {
	// GetUserLoanByUserID returns the cached loans for the given user, stored
	// under key "loan_user_<user_id>". It returns a redis.Nil-wrapped error
	// when the key does not exist.
	GetUserLoanByUserID(ctx context.Context, userID int64) ([]cache.Loan, error)

	// SetUserLoanByUserID caches loans under key "loan_user_<user_id>" as a
	// JSON string, expiring after 10 minutes.
	SetUserLoanByUserID(ctx context.Context, userID int64, loans []cache.Loan) error

	// DeleteUserLoanByUserID removes the cached loans for the given user,
	// stored under key "loan_user_<user_id>".
	DeleteUserLoanByUserID(ctx context.Context, userID int64) error
}

type Service struct {
	db    dbRepoProvider
	cache cacheRepoProvider
}

func NewService(db dbRepoProvider, cache cacheRepoProvider) *Service {
	return &Service{
		db:    db,
		cache: cache,
	}
}
