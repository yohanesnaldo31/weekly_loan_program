package loan

import (
	"context"
	"time"

	"weekly_loan_program/service/loan"
)

type LoanServiceProvider interface {
	// CreateLoanWithBilling inserts loan and its billings within a single
	// transaction, committing only if every insert succeeds.
	CreateLoanWithBilling(ctx context.Context, loan loan.Loan, billings []loan.Billing) (int64, error)

	// GetBillingsByLoanIDAndDueDate returns the billings for the given loan due
	// before dueDate, sorted by due_date ascending. When status is 0, billings
	// are returned regardless of status; otherwise only billings matching status
	// are returned.
	GetBillingsByLoanIDAndDueDate(ctx context.Context, loanID int64, dueDate time.Time, status int16) ([]loan.Billing, error)

	// GetUserLoansByUserID returns up to 10 of the given user's loans, sorted by
	// create_time descending (most recent first). It serves from cache when
	// present, otherwise falls back to the database and populates the cache.
	GetUserLoansByUserID(ctx context.Context, userID int64) ([]loan.Loan, error)

	// UpdateLoanByPayment marks the given billings as paid and updates the
	// loan's total paid amount and status to reflect a payment, all within a
	// single transaction.
	UpdateLoanByPayment(ctx context.Context, input loan.UpdateLoanByPaymentInput) error
}

type Usecase struct {
	loan LoanServiceProvider
}

func NewUsecase(loan LoanServiceProvider) *Usecase {
	return &Usecase{
		loan: loan,
	}
}
