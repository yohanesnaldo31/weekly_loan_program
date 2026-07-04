package db

import "time"

// Loan represents a row in the loan table.
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

// Billing represents a row in the billing table.
type Billing struct {
	ID          int64
	LoanID      int64
	Status      int16
	Amount      int64
	DueDate     time.Time
	PaymentTime *time.Time
}
