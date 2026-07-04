package loan

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"weekly_loan_program/infra/constants"
	"weekly_loan_program/repo/cache"
	"weekly_loan_program/repo/db"
)

// CreateLoanWithBilling inserts loan and its billings within a single
// transaction, committing only if every insert succeeds.
func (s *Service) CreateLoanWithBilling(ctx context.Context, loan Loan, billings []Billing) (int64, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("service: begin transaction for user_id %d: %w", loan.UserID, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	loanID, err := s.db.InsertLoan(ctx, tx, db.Loan(loan))
	if err != nil {
		return 0, fmt.Errorf("service: insert loan for user_id %d: %w", loan.UserID, err)
	}

	for _, billing := range billings {
		billing.LoanID = loanID
		if _, err := s.db.InsertBilling(ctx, tx, db.Billing(billing)); err != nil {
			return 0, fmt.Errorf("service: insert billing for loan_id %d: %w", loanID, err)
		}
	}

	if errCommit := tx.Commit(ctx); errCommit != nil {
		return 0, fmt.Errorf("service: commit transaction for loan_id %d: %w", loanID, err)
	}

	go func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := s.cache.DeleteUserLoanByUserID(cleanupCtx, loan.UserID)
		if err != nil {
			log.Println("error delete redis loan with userID " + strconv.FormatInt(loan.UserID, 10) + ": " + err.Error())
		}
	}()

	return loanID, nil
}

// GetLoansByStatusesAndLastActivityTime returns loans whose status is in the provided list,
// filtered by update_time lower than the provided date and ordered by update_time descending.
func (s *Service) GetLoansByStatusesAndLastActivityTime(ctx context.Context, statuses []int16, lastActivityDate time.Time) ([]Loan, error) {
	loans, err := s.db.GetLoansByStatusesAndLastActivityTime(ctx, statuses, lastActivityDate)
	if err != nil {
		return nil, fmt.Errorf("service: get loans by statuses: %w", err)
	}

	return convertDBLoanToLoan(loans), nil
}

// GetUserLoansByUserID returns up to 10 of the given user's loans, sorted by
// create_time descending (most recent first). It serves from cache when
// present, otherwise falls back to the database and populates the cache.
func (s *Service) GetUserLoansByUserID(ctx context.Context, userID int64) ([]Loan, error) {
	cached, err := s.cache.GetUserLoanByUserID(ctx, userID)
	if err == nil && len(cached) > 0 {
		return convertCacheLoanToLoan(cached), nil
	}

	loans, err := s.db.GetUserLoansByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service: get loans from DB for user_id %d: %w", userID, err)
	}

	if err := s.cache.SetUserLoanByUserID(ctx, userID, convertDBLoanToCacheLoan(loans)); err != nil {
		return nil, fmt.Errorf("service: set cached loans for user_id %d: %w", userID, err)
	}

	return convertDBLoanToLoan(loans), nil
}

// UpdateLoanByPayment marks the given billings as paid and updates the
// loan's total paid amount and status to reflect a payment, all within a
// single transaction.
func (s *Service) UpdateLoanByPayment(ctx context.Context, input UpdateLoanByPaymentInput) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("service: begin transaction for loan_id %d: %w", input.LoanID, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	err = s.db.UpdateBillings(ctx, tx, input.BillingIDs, constants.BILLING_STATUS_PAID, &input.PaymentTime)
	if err != nil {
		return fmt.Errorf("service: error update billings for loan_id %d: %w", input.LoanID, err)
	}

	err = s.db.UpdateLoan(ctx, tx, input.LoanID, input.TotalPaid, int16(input.LoanStatus), input.PaymentTime)
	if err != nil {
		return fmt.Errorf("service: error update loan for loan_id %d: %w", input.LoanID, err)
	}

	if errCommit := tx.Commit(ctx); errCommit != nil {
		return fmt.Errorf("service: commit transaction for loan_id %d: %w", input.LoanID, err)
	}

	go func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := s.cache.DeleteUserLoanByUserID(cleanupCtx, input.UserID)
		if err != nil {
			log.Println("error delete redis loan with userID " + strconv.FormatInt(input.UserID, 10) + ": " + err.Error())
		}
	}()

	return nil
}

// UpdateLoansStatusAndUpdateTimeByIDs updates the status and last activity time
// for all loans whose IDs are included in the provided list.
func (s *Service) UpdateLoansStatusAndUpdateTimeByIDs(ctx context.Context, loanIDs []int64, userIDs []int64, status int16, updateTime time.Time) error {
	if err := s.db.UpdateLoansStatusAndUpdateTimeByIDs(ctx, loanIDs, status, updateTime); err != nil {
		return fmt.Errorf("service: update loans status and update time by ids: %w", err)
	}

	// clean up cache for all userIDs in a separate goroutine to avoid blocking the main flow
	go func(userIDs []int64) {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		for _, userID := range userIDs {
			err := s.cache.DeleteUserLoanByUserID(cleanupCtx, userID)
			if err != nil {
				log.Println("error delete redis loan with userID " + strconv.FormatInt(userID, 10) + ": " + err.Error())
			}
		}
	}(userIDs)

	return nil
}

func convertDBLoanToLoan(loans []db.Loan) []Loan {
	out := make([]Loan, len(loans))
	for i, loan := range loans {
		out[i] = Loan(loan)
	}
	return out
}

func convertCacheLoanToLoan(loans []cache.Loan) []Loan {
	out := make([]Loan, len(loans))
	for i, loan := range loans {
		out[i] = Loan(loan)
	}
	return out
}

func convertDBLoanToCacheLoan(loans []db.Loan) []cache.Loan {
	out := make([]cache.Loan, len(loans))
	for i, loan := range loans {
		out[i] = cache.Loan(loan)
	}
	return out
}
