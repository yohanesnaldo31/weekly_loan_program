package loan

import (
	"context"
	"log"
	"time"
)

type loanUsecaseProvider interface {
	// UpdateLoanDelinquentStatus updates loans to delinquent based on the provided reference time.
	UpdateLoanDelinquentStatus(ctx context.Context, referenceTime time.Time) error
}

type Handler struct {
	loan loanUsecaseProvider
}

func NewHandler(loanUsecase loanUsecaseProvider) *Handler {
	return &Handler{loan: loanUsecase}
}

func (h *Handler) RunDelinquentCheck(ctx context.Context) error {
	log.Println("running delinquent loan check")
	if err := h.loan.UpdateLoanDelinquentStatus(ctx, time.Now()); err != nil {
		log.Printf("delinquent check failed: %v", err)
		return err
	}
	log.Println("delinquent loan check completed")
	return nil
}
