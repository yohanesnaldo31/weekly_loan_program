package loan

import (
	"context"
	"testing"
	"time"
	"weekly_loan_program/service/loan"

	gomock "github.com/golang/mock/gomock"
)

func TestUsecase_RequestLoan(t *testing.T) {
	type mockFields struct {
		loan *MockLoanServiceProvider
	}
	type args struct {
		ctx     context.Context
		request RequestLoanInput
	}
	tests := []struct {
		name    string
		args    args
		mock    func(mf mockFields)
		want    int64
		wantErr bool
	}{
		{
			name: "given_user_has_ongoing_loan",
			args: args{
				ctx: context.Background(),
				request: RequestLoanInput{
					UserID:             1,
					LoanAmount:         1000,
					InstallmentInWeeks: 2,
				},
			},
			mock: func(mf mockFields) {
				mf.loan.EXPECT().GetUserLoansByUserID(context.Background(), int64(1)).Return([]loan.Loan{
					{
						ID:     123,
						Status: 1, // ongoing loan
					},
				}, nil)
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "given_loan_amount_1000_and_weeks_2",
			args: args{
				ctx: context.Background(),
				request: RequestLoanInput{
					UserID:             1,
					LoanAmount:         1000,
					InstallmentInWeeks: 2,
				},
			},
			mock: func(mf mockFields) {
				mf.loan.EXPECT().GetUserLoansByUserID(context.Background(), int64(1)).Return([]loan.Loan{}, nil)
				mf.loan.EXPECT().CreateLoanWithBilling(context.Background(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, loan loan.Loan, billings []loan.Billing) (int64, error) {
					if len(billings) != 2 {
						t.Errorf("expected 2 billings, got %d", len(billings))
					}
					if billings[0].Amount != 550 {
						t.Errorf("expected first billing amount to be 550, got %d", billings[0].Amount)
					}
					if billings[1].Amount != 550 {
						t.Errorf("expected second billing amount to be 550, got %d", billings[1].Amount)
					}

					if billings[0].DueDate.After(billings[1].DueDate) {
						t.Errorf("expected first billing due date to be before second billing due date")
					}

					if billings[0].DueDate.AddDate(0, 0, 7) != billings[1].DueDate {
						t.Errorf("expected second billing due date to be one week after first billing due date")
					}

					if billings[0].Status != 1 {
						t.Errorf("expected first billing status to be 1, got %d", billings[0].Status)
					}
					if billings[1].Status != 1 {
						t.Errorf("expected second billing status to be 1, got %d", billings[1].Status)
					}

					if loan.Status != 1 {
						t.Errorf("expected loan status to be 1, got %d", loan.Status)
					}

					if loan.TotalOutstanding != 1100 {
						t.Errorf("expected loan total outstanding to be 1100, got %d", loan.TotalOutstanding)
					}

					return int64(123), nil
				})
			},
			want:    int64(123),
			wantErr: false,
		},
		{
			name: "given_loan_amount_1000_and_weeks_3",
			args: args{
				ctx: context.Background(),
				request: RequestLoanInput{
					UserID:             1,
					LoanAmount:         1000,
					InstallmentInWeeks: 3,
				},
			},
			mock: func(mf mockFields) {
				mf.loan.EXPECT().GetUserLoansByUserID(context.Background(), int64(1)).Return([]loan.Loan{}, nil)
				mf.loan.EXPECT().CreateLoanWithBilling(context.Background(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, loan loan.Loan, billings []loan.Billing) (int64, error) {
					if len(billings) != 3 {
						t.Errorf("expected 3 billings, got %d", len(billings))
					}
					if billings[0].Amount != 366 {
						t.Errorf("expected first billing amount to be 366, got %d", billings[0].Amount)
					}
					if billings[1].Amount != 366 {
						t.Errorf("expected second billing amount to be 366, got %d", billings[1].Amount)
					}
					if billings[2].Amount != 368 {
						t.Errorf("expected third billing amount to be 368, got %d", billings[2].Amount)
					}

					if billings[0].DueDate.After(billings[1].DueDate) {
						t.Errorf("expected first billing due date to be before second billing due date")
					}

					if billings[0].DueDate.AddDate(0, 0, 7) != billings[1].DueDate {
						t.Errorf("expected second billing due date to be one week after first billing due date")
					}

					if billings[0].Status != 1 {
						t.Errorf("expected first billing status to be 1, got %d", billings[0].Status)
					}
					if billings[1].Status != 1 {
						t.Errorf("expected second billing status to be 1, got %d", billings[1].Status)
					}
					if billings[2].Status != 1 {
						t.Errorf("expected third billing status to be 1, got %d", billings[2].Status)
					}

					if loan.Status != 1 {
						t.Errorf("expected loan status to be 1, got %d", loan.Status)
					}

					if loan.TotalOutstanding != 1100 {
						t.Errorf("expected loan total outstanding to be 1100, got %d", loan.TotalOutstanding)
					}

					return int64(123), nil
				})
			},
			want:    int64(123),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFields := mockFields{
				loan: NewMockLoanServiceProvider(ctrl),
			}
			tt.mock(mockFields)
			uc := &Usecase{
				loan: mockFields.loan,
			}
			got, err := uc.RequestLoan(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.RequestLoan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Usecase.RequestLoan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecase_PayLoan(t *testing.T) {
	mockTime := time.Now()
	type mockFields struct {
		loan *MockLoanServiceProvider
	}
	type args struct {
		ctx     context.Context
		request PayLoanInput
	}
	tests := []struct {
		name    string
		args    args
		mock    func(mf mockFields)
		wantErr bool
	}{
		{
			name: "given_user_has_no_ongoing_loan",
			args: args{
				ctx: context.Background(),
				request: PayLoanInput{
					UserID:        1,
					PaymentAmount: 1000,
				},
			},
			mock: func(mf mockFields) {
				mf.loan.EXPECT().GetUserLoansByUserID(context.Background(), int64(1)).Return([]loan.Loan{}, nil)
			},
			wantErr: true,
		},
		{
			name: "given_user_has_no_ongoing_billings",
			args: args{
				ctx: context.Background(),
				request: PayLoanInput{
					UserID:        1,
					PaymentAmount: 1000,
					PaymentTime:   mockTime,
				},
			},
			mock: func(mf mockFields) {
				mf.loan.EXPECT().GetUserLoansByUserID(context.Background(), int64(1)).Return([]loan.Loan{
					{
						ID:               123,
						Status:           1, // ongoing loan
						TotalOutstanding: 2000,
						TotalPaid:        0,
					},
				}, nil)
				mf.loan.EXPECT().GetBillingsByLoanIDAndDueDate(context.Background(), int64(123), mockTime.AddDate(0, 0, 7), int16(1)).Return([]loan.Billing{}, nil)
			},
			wantErr: true,
		},
		{
			name: "given_user_payment_amount_not_equal_to_total_billings",
			args: args{
				ctx: context.Background(),
				request: PayLoanInput{
					UserID:        1,
					PaymentAmount: 1000,
					PaymentTime:   mockTime,
				},
			},
			mock: func(mf mockFields) {
				mf.loan.EXPECT().GetUserLoansByUserID(context.Background(), int64(1)).Return([]loan.Loan{
					{
						ID:               123,
						Status:           1, // ongoing loan
						TotalOutstanding: 2000,
						TotalPaid:        0,
					},
				}, nil)
				mf.loan.EXPECT().GetBillingsByLoanIDAndDueDate(context.Background(), int64(123), mockTime.AddDate(0, 0, 7), int16(1)).Return([]loan.Billing{
					{
						ID:      456,
						LoanID:  123,
						Amount:  1000,
						DueDate: time.Now(),
						Status:  1, // pending
					},
					{
						ID:      457,
						LoanID:  123,
						Amount:  1000,
						DueDate: time.Now(),
						Status:  1, // pending
					},
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "given_user_has_ongoing_and_stuck_loan",
			args: args{
				ctx: context.Background(),
				request: PayLoanInput{
					UserID:        1,
					PaymentAmount: 2000,
					PaymentTime:   mockTime,
				},
			},
			mock: func(mf mockFields) {
				mf.loan.EXPECT().GetUserLoansByUserID(context.Background(), int64(1)).Return([]loan.Loan{
					{
						ID:               123,
						Status:           1, // ongoing loan
						TotalOutstanding: 2000,
						TotalPaid:        0,
					},
				}, nil)
				mf.loan.EXPECT().GetBillingsByLoanIDAndDueDate(context.Background(), int64(123), mockTime.AddDate(0, 0, 7), int16(1)).Return([]loan.Billing{
					{
						ID:      456,
						LoanID:  123,
						Amount:  1000,
						DueDate: time.Now(),
						Status:  1, // pending
					},
					{
						ID:      457,
						LoanID:  123,
						Amount:  1000,
						DueDate: time.Now(),
						Status:  1, // pending
					},
				}, nil)
				mf.loan.EXPECT().UpdateLoanByPayment(context.Background(), loan.UpdateLoanByPaymentInput{
					UserID:      1,
					LoanID:      123,
					LoanStatus:  4, // complete
					TotalPaid:   2000,
					PaymentTime: mockTime,
					BillingIDs:  []int64{456, 457},
				}).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mf := mockFields{
				loan: NewMockLoanServiceProvider(ctrl),
			}
			tt.mock(mf)
			uc := &Usecase{
				loan: mf.loan,
			}
			if err := uc.PayLoan(tt.args.ctx, tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("Usecase.PayLoan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
