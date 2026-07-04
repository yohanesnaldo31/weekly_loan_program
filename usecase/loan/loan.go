package loan

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"weekly_loan_program/infra/constants"
	"weekly_loan_program/service/loan"
)

const (
	baseDecimal = 10

	interestsRate = 10
)

// RequestLoan creates a new loan for the user along with its weekly billing
// schedule. It rejects the request if the user already has a loan that isn't
// complete. The requested amount is inflated by 10% interest, split evenly
// across the installment weeks, with any rounding leftover added to the last
// billing.
func (uc *Usecase) RequestLoan(ctx context.Context, request RequestLoanInput) (int64, error) {
	userLoans, err := uc.loan.GetUserLoansByUserID(ctx, request.UserID)
	if err != nil {
		log.Println(fmt.Sprintf("error: getting user loans by userID %d: %s", request.UserID, err.Error()))
		return 0, err
	}

	if len(userLoans) != 0 && userLoans[0].Status != constants.LOAN_STATUS_COMPLETE {
		return 0, errors.New("You have ongoing loan, loanID: " + strconv.FormatInt(userLoans[0].ID, baseDecimal))
	}

	// billing calculation
	outstandingAmount := request.LoanAmount*(interestsRate/100) + request.LoanAmount
	billingAmount := outstandingAmount / int64(request.InstallmentInWeeks)
	leftoverAmount := outstandingAmount - billingAmount*int64(request.InstallmentInWeeks)

	currentTime := time.Now()
	currentTimeAfter1Week := currentTime.AddDate(0, 0, 7)
	billingTime := time.Date(currentTimeAfter1Week.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location())

	// billings creation
	billings := make([]loan.Billing, request.InstallmentInWeeks)
	for idx, _ := range billings {
		billings[idx] = loan.Billing{
			Status:  constants.BILLING_STATUS_IN_PROGRESS,
			Amount:  billingAmount,
			DueDate: billingTime,
		}
		billingTime = billingTime.AddDate(0, 0, 7)
	}

	if leftoverAmount > 0 {
		billings[len(billings)-1].Amount += leftoverAmount
	}

	// loan creation
	loan := loan.Loan{
		UserID:           request.UserID,
		Status:           constants.LOAN_STATUS_NEW,
		LoanAmount:       request.LoanAmount,
		TotalOutstanding: outstandingAmount,
		TotalWeek:        int16(request.InstallmentInWeeks),
		TotalPaid:        0,
		CreateTime:       currentTime,
	}

	loanID, err := uc.loan.CreateLoanWithBilling(ctx, loan, billings)
	if err != nil {
		log.Println(fmt.Sprintf("error: creating loan by userID %d: %s", request.UserID, err.Error()))
		return 0, err
	}
	return loanID, nil
}

// PayLoan applies a payment towards the user's ongoing loan. It rejects the
// request if the user has no ongoing loan, has no billing due within the
// next 7 days, or if the payment amount doesn't exactly match the total of
// those due billings. On success, it marks the due billings as paid and
// updates the loan's total paid and status (completing the loan if the
// payment covers the remaining outstanding amount).
func (uc *Usecase) PayLoan(ctx context.Context, request PayLoanInput) error {
	userLoans, err := uc.loan.GetUserLoansByUserID(ctx, request.UserID)
	if err != nil {
		log.Println(fmt.Sprintf("error: getting user loans by userID %d: %s", request.UserID, err.Error()))
		return err
	}

	if len(userLoans) == 0 || userLoans[0].Status == constants.LOAN_STATUS_COMPLETE {
		return errors.New("You have no ongoing loan")
	}

	currentTime := request.PaymentTime
	if currentTime.IsZero() {
		currentTime = time.Now()
	}

	userLoan := userLoans[0]
	dueDateLimit := currentTime.AddDate(0, 0, 7)
	billings, err := uc.loan.GetBillingsByLoanIDAndDueDate(ctx, userLoan.ID, dueDateLimit, constants.BILLING_STATUS_IN_PROGRESS)
	if err != nil {
		log.Println(fmt.Sprintf("error: getting loan bills by loanID %d: %s", userLoan.ID, err.Error()))
		return err
	}

	if len(billings) == 0 {
		return errors.New("You have no ongoing bill or already paid this week billing")
	}

	var totalBillings int64
	billingsID := make([]int64, 0)
	for _, billing := range billings {
		billingsID = append(billingsID, billing.ID)
		totalBillings += billing.Amount
	}

	// payment have to be equal to available billings
	if totalBillings != request.PaymentAmount {
		return errors.New("You have to pay for this amount: " + strconv.FormatInt(totalBillings, 10))
	}

	// validate if user already finish paying the loan
	loanStatus := constants.LOAN_STATUS_IN_PROGRESS
	if userLoan.TotalOutstanding <= userLoan.TotalPaid+request.PaymentAmount {
		loanStatus = constants.LOAN_STATUS_COMPLETE
	}

	err = uc.loan.UpdateLoanByPayment(ctx, loan.UpdateLoanByPaymentInput{
		UserID:      request.UserID,
		LoanID:      userLoan.ID,
		LoanStatus:  loanStatus,
		TotalPaid:   userLoan.TotalPaid + request.PaymentAmount,
		PaymentTime: currentTime,
		BillingIDs:  billingsID,
	})
	if err != nil {
		log.Println(fmt.Sprintf("error: do payment for uesr %d: %s", request.UserID, err.Error()))
		return err
	}

	return nil
}

// GetUserLoansByUserID returns up to 10 of the given user's loans, sorted by
// create_time descending (most recent first).
func (uc *Usecase) GetUserLoansByUserID(ctx context.Context, userID int64) ([]Loan, error) {
	svcLoans, err := uc.loan.GetUserLoansByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	loans := make([]Loan, len(svcLoans))
	for i, svcLoan := range svcLoans {
		loans[i] = Loan(svcLoan)
	}
	return loans, nil
}
