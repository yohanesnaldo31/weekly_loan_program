package loan

import "time"

// LoanResponse is the JSON representation of a single loan.
type LoanResponse struct {
	ID                 int64     `json:"id"`
	UserID             int64     `json:"user_id"`
	Status             int16     `json:"status"`
	StatusName         string    `json:"status_name"`
	LoanAmount         int64     `json:"loan_amount"`
	CurrentOutstanding int64     `json:"current_outstanding"`
	TotalOutstanding   int64     `json:"total_outstanding"`
	TotalPaid          int64     `json:"total_paid"`
	TotalWeek          int16     `json:"total_week"`
	CreateTime         time.Time `json:"create_time"`
	UpdateTime         time.Time `json:"update_time,omitempty"`
}

// GetUserLoansResponse is the JSON body returned by GetUserLoans.
type GetUserLoansResponse struct {
	Loans        []LoanResponse `json:"loans,omitempty"`
	IsSuccess    bool           `json:"is_success"`
	ErrorMessage string         `json:"error_message,omitempty"`
}

// CheckUserDelinquentResponse is the JSON body returned by CheckUserDelinquent.
type CheckUserDelinquentResponse struct {
	IsDelinquent bool   `json:"is_delinquent"`
	IsSuccess    bool   `json:"is_success"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// GetUserOutstandingResponse is the JSON body returned by GetUserOutstandingHandler.
type GetUserOutstandingResponse struct {
	OutstandingAmount int64  `json:"outstanding_amount"`
	IsSuccess         bool   `json:"is_success"`
	ErrorMessage      string `json:"error_message,omitempty"`
}

// RequestLoanRequest is the JSON body accepted by RequestLoan.
type RequestLoanRequest struct {
	UserID             int64 `json:"user_id"`
	LoanAmount         int64 `json:"loan_amount"`
	InstallmentInWeeks int32 `json:"installment_in_weeks"`
}

// RequestLoanResponse is the JSON body returned by RequestLoan.
type RequestLoanResponse struct {
	LoanID       int64  `json:"loan_id"`
	IsSuccess    bool   `json:"is_success"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// PayLoanRequest is the JSON body accepted by PayLoan.
type PayLoanRequest struct {
	UserID        int64  `json:"user_id"`
	PaymentAmount int64  `json:"payment_amount"`
	PaymentTime   string `json:"payment_time,omitempty"` // format: YYYY-MM-DD, optional
}

// PayLoanResponse is the JSON body returned by PayLoan.
type PayLoanResponse struct {
	IsSuccess    bool   `json:"is_success"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// TriggerDelinquentCheckRequest is the JSON body accepted by TriggerDelinquentCheck.
type TriggerDelinquentCheckRequest struct {
	Time string `json:"time"`
}

// TriggerDelinquentCheckResponse is the JSON body returned by TriggerDelinquentCheck.
type TriggerDelinquentCheckResponse struct {
	IsSuccess    bool   `json:"is_success"`
	ErrorMessage string `json:"error_message,omitempty"`
}
