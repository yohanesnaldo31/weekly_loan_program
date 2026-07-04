package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	queryInsertBilling = `
		INSERT INTO billing (loan_id, status, amount, due_date, payment_time)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	queryGetBillingsByLoanID = `
		SELECT id, loan_id, status, amount, due_date, payment_time
		FROM billing
		WHERE loan_id = $1 AND due_date < $2 AND ($3 = 0 OR status = $3)
		ORDER BY due_date ASC
	`

	queryUpdateBillings = `
		UPDATE billing
		SET status = $1, payment_time = $2
		WHERE id = ANY($3)
	`
)

// InsertBilling inserts billing within the given transaction and returns the
// generated billing ID.
func (r *Repository) InsertBilling(ctx context.Context, tx pgx.Tx, billing Billing) (int64, error) {
	var billingID int64
	err := tx.QueryRow(ctx, queryInsertBilling,
		billing.LoanID,
		billing.Status,
		billing.Amount,
		billing.DueDate,
		billing.PaymentTime,
	).Scan(&billingID)
	if err != nil {
		return 0, fmt.Errorf("db: insert billing for loan_id %d: %w", billing.LoanID, err)
	}

	return billingID, nil
}

// GetBillingsByLoanIDAndDueDate returns the billings for the given loan due
// before dueDate, sorted by due_date ascending. When status is 0, billings
// are returned regardless of status; otherwise only billings matching status
// are returned.
func (r *Repository) GetBillingsByLoanIDAndDueDate(ctx context.Context, loanID int64, dueDate time.Time, status int16) ([]Billing, error) {
	rows, err := r.db.Query(ctx, queryGetBillingsByLoanID, loanID, dueDate, status)
	if err != nil {
		return nil, fmt.Errorf("db: query billings by loan id: %w", err)
	}
	defer rows.Close()

	var billings []Billing
	for rows.Next() {
		var billing Billing
		if err := rows.Scan(
			&billing.ID,
			&billing.LoanID,
			&billing.Status,
			&billing.Amount,
			&billing.DueDate,
			&billing.PaymentTime,
		); err != nil {
			return nil, fmt.Errorf("db: scan billing row for loan_id %d: %w", loanID, err)
		}
		billings = append(billings, billing)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db: iterate billing rows for loan_id %d: %w", loanID, err)
	}

	return billings, nil
}

// UpdateBillings updates the status and payment time of the billings with
// the given IDs within the given transaction.
func (r *Repository) UpdateBillings(ctx context.Context, tx pgx.Tx, billingIDs []int64, status int16, paymentTime *time.Time) error {
	if _, err := tx.Exec(ctx, queryUpdateBillings, status, paymentTime, billingIDs); err != nil {
		return fmt.Errorf("db: update billing_ids %v: %w", billingIDs, err)
	}

	return nil
}
