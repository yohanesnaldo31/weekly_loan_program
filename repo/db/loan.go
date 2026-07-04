package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	queryGetUserLoansByUserID = `
		SELECT id, user_id, status, loan_amount, total_outstanding, total_paid, total_week, create_time, update_time
		FROM loan
		WHERE user_id = $1
		ORDER BY create_time DESC
		LIMIT 10
	`

	queryGetLoanByID = `
		SELECT id, user_id, status, loan_amount, total_outstanding, total_paid, total_week, create_time, update_time
		FROM loan
		WHERE id = $1
	`

	queryInsertLoan = `
		INSERT INTO loan (user_id, status, loan_amount, total_outstanding, total_paid, total_week, create_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	queryUpdateLoan = `
		UPDATE loan
		SET total_paid = $1, status = $2, update_time = $3
		WHERE id = $4
	`
)

// GetUserLoansByUserID returns up to 10 of the given user's loans, sorted by
// create_time descending (most recent first).
func (r *Repository) GetUserLoansByUserID(ctx context.Context, userID int64) ([]Loan, error) {
	rows, err := r.db.Query(ctx, queryGetUserLoansByUserID, userID)
	if err != nil {
		return nil, fmt.Errorf("db: query loans by user id: %w", err)
	}
	defer rows.Close()

	var loans []Loan
	for rows.Next() {
		var loan Loan
		if err := rows.Scan(
			&loan.ID,
			&loan.UserID,
			&loan.Status,
			&loan.LoanAmount,
			&loan.TotalOutstanding,
			&loan.TotalPaid,
			&loan.TotalWeek,
			&loan.CreateTime,
			&loan.UpdateTime,
		); err != nil {
			return nil, fmt.Errorf("db: scan loan row for user_id %d: %w", userID, err)
		}
		loans = append(loans, loan)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db: iterate loan rows for user_id %d: %w", userID, err)
	}

	return loans, nil
}

// InsertLoan inserts loan within the given transaction and returns the
// generated loan ID.
func (r *Repository) InsertLoan(ctx context.Context, tx pgx.Tx, loan Loan) (int64, error) {
	var loanID int64
	err := tx.QueryRow(ctx, queryInsertLoan,
		loan.UserID,
		loan.Status,
		loan.LoanAmount,
		loan.TotalOutstanding,
		loan.TotalPaid,
		loan.TotalWeek,
		loan.CreateTime,
	).Scan(&loanID)
	if err != nil {
		return 0, fmt.Errorf("db: insert loan for user_id %d: %w", loan.UserID, err)
	}

	return loanID, nil
}

// UpdateLoan updates the total paid, status, and update time of the loan
// with the given ID within the given transaction.
func (r *Repository) UpdateLoan(ctx context.Context, tx pgx.Tx, loanID int64, totalPaid int64, status int16, updateTime time.Time) error {
	if _, err := tx.Exec(ctx, queryUpdateLoan, totalPaid, status, updateTime, loanID); err != nil {
		return fmt.Errorf("db: update loan_id %d: %w", loanID, err)
	}

	return nil
}
