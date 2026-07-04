package loan

import (
	"context"
	"fmt"
	"time"
)

// GetBillingsByLoanIDAndDueDate returns the billings for the given loan due
// before dueDate, sorted by due_date ascending. When status is 0, billings
// are returned regardless of status; otherwise only billings matching status
// are returned.
func (s *Service) GetBillingsByLoanIDAndDueDate(ctx context.Context, loanID int64, dueDate time.Time, status int16) ([]Billing, error) {
	billings, err := s.db.GetBillingsByLoanIDAndDueDate(ctx, loanID, dueDate, status)
	if err != nil {
		return nil, fmt.Errorf("service: get billings for loan_id %d: %w", loanID, err)
	}

	result := make([]Billing, len(billings))
	for i, b := range billings {
		result[i] = Billing(b)
	}
	return result, nil
}
