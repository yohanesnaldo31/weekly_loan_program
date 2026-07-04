package loan

import (
	"context"

	"weekly_loan_program/usecase/loan"
)

type loanUsecaseProvider interface {
	// RequestLoan creates a new loan for the user along with its weekly billing
	// schedule. It rejects the request if the user already has a loan that isn't
	// complete. The requested amount is inflated by 10% interest, split evenly
	// across the installment weeks, with any rounding leftover added to the last
	// billing.
	RequestLoan(ctx context.Context, request loan.RequestLoanInput) (int64, error)

	// GetUserLoansByUserID returns up to 10 of the given user's loans, sorted by
	// create_time descending (most recent first).
	GetUserLoansByUserID(ctx context.Context, userID int64) ([]loan.Loan, error)

	// PayLoan applies a payment towards the user's ongoing loan. It rejects the
	// request if the user has no ongoing loan, has no billing due within the
	// next 7 days, or if the payment amount doesn't exactly match the total of
	// those due billings. On success, it marks the due billings as paid and
	// updates the loan's total paid and status (completing the loan if the
	// payment covers the remaining outstanding amount).
	PayLoan(ctx context.Context, request loan.PayLoanInput) error
}

type Handler struct {
	loan loanUsecaseProvider
}

func NewHandler(loan loanUsecaseProvider) *Handler {
	return &Handler{
		loan: loan,
	}
}
