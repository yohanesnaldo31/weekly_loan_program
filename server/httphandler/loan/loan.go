package loan

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"weekly_loan_program/server/httphandler"
	loanUC "weekly_loan_program/usecase/loan"
)

const dateLayout = "2006-01-02"

// GetUserLoansHandler handles GET requests that return a user's loans.
func (h *Handler) GetUserLoansHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	if err != nil || userID <= 0 {
		httphandler.WriteRequestResponse(w, http.StatusBadRequest, GetUserLoansResponse{
			IsSuccess:    false,
			ErrorMessage: "invalid user_id",
		})
		return
	}

	loans, err := h.loan.GetUserLoansByUserID(r.Context(), userID)
	if err != nil {
		httphandler.WriteRequestResponse(w, http.StatusInternalServerError, GetUserLoansResponse{
			IsSuccess:    false,
			ErrorMessage: err.Error(),
		})
		return
	}

	httphandler.WriteRequestResponse(w, http.StatusOK, GetUserLoansResponse{
		Loans:     convertToLoanResponse(loans),
		IsSuccess: true,
	})
}

// RequestLoanHandler handles POST requests that create a new loan for a user.
func (h *Handler) RequestLoanHandler(w http.ResponseWriter, r *http.Request) {
	var req RequestLoanRequest

	// request validation
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httphandler.WriteRequestResponse(w, http.StatusBadRequest, RequestLoanResponse{
			IsSuccess:    false,
			ErrorMessage: "invalid request body",
		})
		return
	}

	if req.UserID <= 0 {
		httphandler.WriteRequestResponse(w, http.StatusBadRequest, RequestLoanResponse{
			IsSuccess:    false,
			ErrorMessage: "invalid user_id",
		})
		return
	}

	if req.LoanAmount <= 10000 { // not sure what the requirements of the loan minimum or maximum
		httphandler.WriteRequestResponse(w, http.StatusBadRequest, RequestLoanResponse{
			IsSuccess:    false,
			ErrorMessage: "invalid user_id",
		})
		return
	}

	if req.InstallmentInWeeks <= 0 ||
		int64(req.InstallmentInWeeks) > req.LoanAmount { // shouldn't be lower than the loan amount
		httphandler.WriteRequestResponse(w, http.StatusBadRequest, RequestLoanResponse{
			IsSuccess:    false,
			ErrorMessage: "invalid installment_in_weeks",
		})
		return
	}

	if req.InstallmentInWeeks > 261 {
		httphandler.WriteRequestResponse(w, http.StatusBadRequest, RequestLoanResponse{
			IsSuccess:    false,
			ErrorMessage: "installment can't be more than 5 years",
		})
		return
	}

	// calling business logic
	loanID, err := h.loan.RequestLoan(r.Context(), loanUC.RequestLoanInput{
		UserID:             req.UserID,
		LoanAmount:         req.LoanAmount,
		InstallmentInWeeks: req.InstallmentInWeeks,
	})
	if err != nil {
		httphandler.WriteRequestResponse(w, http.StatusInternalServerError, RequestLoanResponse{
			IsSuccess:    false,
			ErrorMessage: err.Error(),
		})
		return
	}

	httphandler.WriteRequestResponse(w, http.StatusOK, RequestLoanResponse{
		LoanID:    loanID,
		IsSuccess: true,
	})
}

// PayLoanHandler handles POST requests that apply a payment to a user's
// ongoing loan.
func (h *Handler) PayLoanHandler(w http.ResponseWriter, r *http.Request) {
	var req PayLoanRequest

	// request validation
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httphandler.WriteRequestResponse(w, http.StatusBadRequest, PayLoanResponse{
			IsSuccess:    false,
			ErrorMessage: "invalid request body",
		})
		return
	}

	if req.UserID <= 0 {
		httphandler.WriteRequestResponse(w, http.StatusBadRequest, PayLoanResponse{
			IsSuccess:    false,
			ErrorMessage: "invalid user_id",
		})
		return
	}

	if req.PaymentAmount <= 0 {
		httphandler.WriteRequestResponse(w, http.StatusBadRequest, PayLoanResponse{
			IsSuccess:    false,
			ErrorMessage: "invalid payment_amount",
		})
		return
	}

	var paymentTime time.Time
	if req.PaymentTime != "" {
		var err error
		paymentTime, err = time.Parse(dateLayout, req.PaymentTime)
		if err != nil {
			httphandler.WriteRequestResponse(w, http.StatusBadRequest, PayLoanResponse{
				IsSuccess:    false,
				ErrorMessage: "invalid payment_time, expected format YYYY-MM-DD",
			})
			return
		}
	}

	// calling business logic
	if err := h.loan.PayLoan(r.Context(), loanUC.PayLoanInput{
		UserID:        req.UserID,
		PaymentAmount: req.PaymentAmount,
		PaymentTime:   paymentTime,
	}); err != nil {
		httphandler.WriteRequestResponse(w, http.StatusInternalServerError, PayLoanResponse{
			IsSuccess:    false,
			ErrorMessage: err.Error(),
		})
		return
	}

	httphandler.WriteRequestResponse(w, http.StatusOK, PayLoanResponse{
		IsSuccess: true,
	})
}

func convertToLoanResponse(loans []loanUC.Loan) []LoanResponse {
	out := make([]LoanResponse, len(loans))
	for i, loan := range loans {
		out[i] = LoanResponse(loan)
	}
	return out
}
