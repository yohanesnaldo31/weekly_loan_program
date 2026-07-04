package loan

import "time"

type RequestLoanInput struct {
	UserID             int64
	LoanAmount         int64
	InstallmentInWeeks int32
}

type PayLoanInput struct {
	UserID        int64
	PaymentAmount int64
	PaymentTime   time.Time
}

// Loan represents a loan returned by the usecase layer.
type Loan struct {
	ID               int64
	UserID           int64
	Status           int16
	LoanAmount       int64
	TotalOutstanding int64
	TotalPaid        int64
	TotalWeek        int16
	CreateTime       time.Time
	UpdateTime       *time.Time
}
