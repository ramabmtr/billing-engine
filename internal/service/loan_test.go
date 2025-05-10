package service

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github/ramabmtr/billing-engine/internal/constant"
	"github/ramabmtr/billing-engine/internal/model"
	"github/ramabmtr/billing-engine/internal/repository"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockLoanRepo is a mock implementation of repository.LoanRepo
type MockLoanRepo struct {
	mock.Mock
}

func (m *MockLoanRepo) WithTx(tx *gorm.DB) repository.LoanRepo {
	args := m.Called(tx)
	return args.Get(0).(repository.LoanRepo)
}

func (m *MockLoanRepo) Create(l *model.Loan) error {
	args := m.Called(l)
	return args.Error(0)
}

func (m *MockLoanRepo) Get(l *model.Loan) error {
	args := m.Called(l)
	// Simulate the behavior of Get by setting fields on the loan
	if args.Error(0) == nil && l != nil {
		l.Principal = decimal.NewFromInt(5_000_000)
		l.AnnualInterestRate = decimal.NewFromInt(10)
		l.Period = 50
		l.PeriodUnit = constant.PeriodUnitWeek
		l.TotalRepayment = decimal.NewFromInt(5_500_000)
		l.CreatedAt = time.Now().UTC()
	}
	return args.Error(0)
}

func (m *MockLoanRepo) FindByBorrowerID(borrowerID string) ([]*model.LoanWithCompleteStatus, error) {
	args := m.Called(borrowerID)
	return args.Get(0).([]*model.LoanWithCompleteStatus), args.Error(1)
}

// MockLoanPaymentRepo is a mock implementation of repository.LoanPaymentRepo
type MockLoanPaymentRepo struct {
	mock.Mock
}

func (m *MockLoanPaymentRepo) WithTx(tx *gorm.DB) repository.LoanPaymentRepo {
	args := m.Called(tx)
	return args.Get(0).(repository.LoanPaymentRepo)
}

func (m *MockLoanPaymentRepo) CreateBulk(lps []*model.LoanPayment) error {
	args := m.Called(lps)
	return args.Error(0)
}

func (m *MockLoanPaymentRepo) GetTotalOutstandingByLoanID(loanID string) (decimal.Decimal, error) {
	args := m.Called(loanID)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

func (m *MockLoanPaymentRepo) GetTotalOutstandingByBorrowerID(borrowerID string) (decimal.Decimal, error) {
	args := m.Called(borrowerID)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

func (m *MockLoanPaymentRepo) Find(lp model.LoanPayment) ([]*model.LoanPayment, error) {
	args := m.Called(lp)
	return args.Get(0).([]*model.LoanPayment), args.Error(1)
}

func (m *MockLoanPaymentRepo) ChangeStatusToPaid(loanIds []string, paidAt time.Time) error {
	args := m.Called(loanIds, paidAt)
	return args.Error(0)
}

// MockLockManager is a mock implementation of the LockManager interface used in LoanService
type MockLockManager struct {
	mock.Mock
}

func (m *MockLockManager) GetLock(key string) *sync.Mutex {
	args := m.Called(key)
	return args.Get(0).(*sync.Mutex)
}

func TestLoanService_CreateLoanRequest(t *testing.T) {
	tests := []struct {
		name          string
		borrowerID    string
		mockSetup     func(mockLoanRepo *MockLoanRepo, mockLoanPaymentRepo *MockLoanPaymentRepo)
		expectedError bool
	}{
		{
			name:       "Success",
			borrowerID: "borrower-id-1",
			mockSetup: func(mockLoanRepo *MockLoanRepo, mockLoanPaymentRepo *MockLoanPaymentRepo) {
				// No outstanding amount
				mockLoanPaymentRepo.On("GetTotalOutstandingByBorrowerID", "borrower-id-1").
					Return(decimal.NewFromInt(0), nil)

				// Transaction handling
				mockLoanRepo.On("WithTx", mock.Anything).Return(mockLoanRepo)
				mockLoanPaymentRepo.On("WithTx", mock.Anything).Return(mockLoanPaymentRepo)

				// Create loan
				mockLoanRepo.On("Create", mock.MatchedBy(func(l *model.Loan) bool {
					return l.BorrowerID == "borrower-id-1" &&
						l.Principal.Equal(decimal.NewFromInt(5_000_000)) &&
						l.AnnualInterestRate.Equal(decimal.NewFromInt(10)) &&
						l.Period == 50 &&
						l.PeriodUnit == constant.PeriodUnitWeek
				})).Return(nil)

				// Create loan payments
				mockLoanPaymentRepo.On("CreateBulk", mock.Anything).Return(nil)
			},
			expectedError: false,
		},
		{
			name:       "Outstanding Amount Exists",
			borrowerID: "borrower-id-2",
			mockSetup: func(mockLoanRepo *MockLoanRepo, mockLoanPaymentRepo *MockLoanPaymentRepo) {
				// Outstanding amount exists
				mockLoanPaymentRepo.On("GetTotalOutstandingByBorrowerID", "borrower-id-2").
					Return(decimal.NewFromInt(1000), nil)
			},
			expectedError: true,
		},
		{
			name:       "Error Getting Outstanding Amount",
			borrowerID: "borrower-id-3",
			mockSetup: func(mockLoanRepo *MockLoanRepo, mockLoanPaymentRepo *MockLoanPaymentRepo) {
				// Error getting outstanding amount
				mockLoanPaymentRepo.On("GetTotalOutstandingByBorrowerID", "borrower-id-3").
					Return(decimal.NewFromInt(0), errors.New("database error"))
			},
			expectedError: true,
		},
		{
			name:       "Error Creating Loan",
			borrowerID: "borrower-id-4",
			mockSetup: func(mockLoanRepo *MockLoanRepo, mockLoanPaymentRepo *MockLoanPaymentRepo) {
				// No outstanding amount
				mockLoanPaymentRepo.On("GetTotalOutstandingByBorrowerID", "borrower-id-4").
					Return(decimal.NewFromInt(0), nil)

				// Transaction handling
				mockLoanRepo.On("WithTx", mock.Anything).Return(mockLoanRepo)

				// Error creating loan
				mockLoanRepo.On("Create", mock.Anything).Return(errors.New("database error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLoanRepo := new(MockLoanRepo)
			mockLoanPaymentRepo := new(MockLoanPaymentRepo)
			tt.mockSetup(mockLoanRepo, mockLoanPaymentRepo)

			service := NewLoanService(mockLoanRepo, mockLoanPaymentRepo)
			loan, err := service.CreateLoanRequest(tt.borrowerID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, loan)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, loan)
				assert.Equal(t, tt.borrowerID, loan.BorrowerID)
				assert.Equal(t, decimal.NewFromInt(5_000_000), loan.Principal)
				assert.Equal(t, decimal.NewFromInt(10), loan.AnnualInterestRate)
				assert.Equal(t, 50, loan.Period)
				assert.Equal(t, constant.PeriodUnitWeek, string(loan.PeriodUnit))
				assert.NotEmpty(t, loan.ID)
			}

			mockLoanRepo.AssertExpectations(t)
			mockLoanPaymentRepo.AssertExpectations(t)
		})
	}
}

func TestLoanService_GetLoansByBorrowerID(t *testing.T) {
	tests := []struct {
		name          string
		borrowerID    string
		mockSetup     func(mockLoanRepo *MockLoanRepo)
		expectedError bool
		expectedCount int
	}{
		{
			name:       "Success with loans",
			borrowerID: "borrower-id-1",
			mockSetup: func(mockLoanRepo *MockLoanRepo) {
				loans := []*model.LoanWithCompleteStatus{
					{
						Loan: model.Loan{
							ID:                 uuid.Must(uuid.NewV7()).String(),
							BorrowerID:         "borrower-id-1",
							Principal:          decimal.NewFromInt(5_000_000),
							AnnualInterestRate: decimal.NewFromInt(10),
							Period:             50,
							PeriodUnit:         constant.PeriodUnitWeek,
							TotalRepayment:     decimal.NewFromInt(5_500_000),
							CreatedAt:          time.Now().UTC(),
						},
						IsCompleted: false,
					},
					{
						Loan: model.Loan{
							ID:                 uuid.Must(uuid.NewV7()).String(),
							BorrowerID:         "borrower-id-1",
							Principal:          decimal.NewFromInt(3_000_000),
							AnnualInterestRate: decimal.NewFromInt(10),
							Period:             30,
							PeriodUnit:         constant.PeriodUnitWeek,
							TotalRepayment:     decimal.NewFromInt(3_300_000),
							CreatedAt:          time.Now().UTC(),
						},
						IsCompleted: true,
					},
				}
				mockLoanRepo.On("FindByBorrowerID", "borrower-id-1").Return(loans, nil)
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name:       "Success with empty list",
			borrowerID: "borrower-id-2",
			mockSetup: func(mockLoanRepo *MockLoanRepo) {
				loans := []*model.LoanWithCompleteStatus{}
				mockLoanRepo.On("FindByBorrowerID", "borrower-id-2").Return(loans, nil)
			},
			expectedError: false,
			expectedCount: 0,
		},
		{
			name:       "Repository Error",
			borrowerID: "borrower-id-3",
			mockSetup: func(mockLoanRepo *MockLoanRepo) {
				mockLoanRepo.On("FindByBorrowerID", "borrower-id-3").
					Return([]*model.LoanWithCompleteStatus{}, errors.New("database error"))
			},
			expectedError: true,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLoanRepo := new(MockLoanRepo)
			mockLoanPaymentRepo := new(MockLoanPaymentRepo)
			tt.mockSetup(mockLoanRepo)

			service := NewLoanService(mockLoanRepo, mockLoanPaymentRepo)
			loans, err := service.GetLoansByBorrowerID(tt.borrowerID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, loans, tt.expectedCount)
				if tt.expectedCount > 0 {
					assert.Equal(t, tt.borrowerID, loans[0].BorrowerID)
				}
			}

			mockLoanRepo.AssertExpectations(t)
		})
	}
}

func TestLoanService_GetLoanDetail(t *testing.T) {
	tests := []struct {
		name          string
		loanID        string
		mockSetup     func(mockLoanRepo *MockLoanRepo, mockLoanPaymentRepo *MockLoanPaymentRepo)
		expectedError bool
	}{
		{
			name:   "Success",
			loanID: "loan-id-1",
			mockSetup: func(mockLoanRepo *MockLoanRepo, mockLoanPaymentRepo *MockLoanPaymentRepo) {
				mockLoanRepo.On("Get", mock.MatchedBy(func(l *model.Loan) bool {
					return l.ID == "loan-id-1"
				})).Return(nil)
				mockLoanPaymentRepo.On("GetTotalOutstandingByLoanID", "loan-id-1").
					Return(decimal.NewFromInt(2_000_000), nil)
			},
			expectedError: false,
		},
		{
			name:   "Error Getting Loan",
			loanID: "loan-id-2",
			mockSetup: func(mockLoanRepo *MockLoanRepo, mockLoanPaymentRepo *MockLoanPaymentRepo) {
				mockLoanRepo.On("Get", mock.MatchedBy(func(l *model.Loan) bool {
					return l.ID == "loan-id-2"
				})).Return(errors.New("database error"))
			},
			expectedError: true,
		},
		{
			name:   "Error Getting Outstanding Amount",
			loanID: "loan-id-3",
			mockSetup: func(mockLoanRepo *MockLoanRepo, mockLoanPaymentRepo *MockLoanPaymentRepo) {
				mockLoanRepo.On("Get", mock.MatchedBy(func(l *model.Loan) bool {
					return l.ID == "loan-id-3"
				})).Return(nil)
				mockLoanPaymentRepo.On("GetTotalOutstandingByLoanID", "loan-id-3").
					Return(decimal.NewFromInt(0), errors.New("database error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLoanRepo := new(MockLoanRepo)
			mockLoanPaymentRepo := new(MockLoanPaymentRepo)
			tt.mockSetup(mockLoanRepo, mockLoanPaymentRepo)

			service := NewLoanService(mockLoanRepo, mockLoanPaymentRepo)
			loan, outstanding, err := service.GetLoanDetail(tt.loanID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, loan)
				assert.Equal(t, tt.loanID, loan.ID)
				assert.Equal(t, decimal.NewFromInt(2_000_000), outstanding)
			}

			mockLoanRepo.AssertExpectations(t)
			mockLoanPaymentRepo.AssertExpectations(t)
		})
	}
}

func TestLoanService_GetLoanPaymentsByLoanID(t *testing.T) {
	tests := []struct {
		name          string
		loanID        string
		mockSetup     func(mockLoanPaymentRepo *MockLoanPaymentRepo)
		expectedError bool
		expectedCount int
	}{
		{
			name:   "Success with loan payments",
			loanID: "loan-id-1",
			mockSetup: func(mockLoanPaymentRepo *MockLoanPaymentRepo) {
				loanPayments := []*model.LoanPayment{
					{
						ID:         uuid.Must(uuid.NewV7()).String(),
						LoanID:     "loan-id-1",
						BorrowerID: "borrower-id-1",
						Amount:     decimal.NewFromInt(110_000),
						DueDate:    time.Now().UTC().AddDate(0, 0, -7),
						Status:     constant.LoanPaymentStatusPaid,
						PaidAt:     func() *time.Time { t := time.Now().UTC().AddDate(0, 0, -5); return &t }(),
					},
					{
						ID:         uuid.Must(uuid.NewV7()).String(),
						LoanID:     "loan-id-1",
						BorrowerID: "borrower-id-1",
						Amount:     decimal.NewFromInt(110_000),
						DueDate:    time.Now().UTC().AddDate(0, 0, 7),
						Status:     constant.LoanPaymentStatusUnpaid,
					},
				}
				mockLoanPaymentRepo.On("Find", mock.MatchedBy(func(lp model.LoanPayment) bool {
					return lp.LoanID == "loan-id-1"
				})).Return(loanPayments, nil)
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name:   "Success with empty list",
			loanID: "loan-id-2",
			mockSetup: func(mockLoanPaymentRepo *MockLoanPaymentRepo) {
				loanPayments := []*model.LoanPayment{}
				mockLoanPaymentRepo.On("Find", mock.MatchedBy(func(lp model.LoanPayment) bool {
					return lp.LoanID == "loan-id-2"
				})).Return(loanPayments, nil)
			},
			expectedError: false,
			expectedCount: 0,
		},
		{
			name:   "Repository Error",
			loanID: "loan-id-3",
			mockSetup: func(mockLoanPaymentRepo *MockLoanPaymentRepo) {
				mockLoanPaymentRepo.On("Find", mock.MatchedBy(func(lp model.LoanPayment) bool {
					return lp.LoanID == "loan-id-3"
				})).Return([]*model.LoanPayment{}, errors.New("database error"))
			},
			expectedError: true,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLoanRepo := new(MockLoanRepo)
			mockLoanPaymentRepo := new(MockLoanPaymentRepo)
			tt.mockSetup(mockLoanPaymentRepo)

			service := NewLoanService(mockLoanRepo, mockLoanPaymentRepo)
			loanPayments, err := service.GetLoanPaymentsByLoanID(tt.loanID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, loanPayments, tt.expectedCount)
				if tt.expectedCount > 0 {
					assert.Equal(t, tt.loanID, loanPayments[0].LoanID)
				}
			}

			mockLoanPaymentRepo.AssertExpectations(t)
		})
	}
}

func TestLoanService_MakePayment(t *testing.T) {
	now := time.Now().UTC()
	pastDue := now.AddDate(0, 0, -7)
	futureDue := now.AddDate(0, 0, 7)

	tests := []struct {
		name          string
		loanID        string
		amount        decimal.Decimal
		mockSetup     func(mockLoanPaymentRepo *MockLoanPaymentRepo, mockLockManager *MockLockManager)
		expectedError bool
	}{
		{
			name:   "Success - Pay Exact Amount",
			loanID: "loan-id-1",
			amount: decimal.NewFromInt(110_000),
			mockSetup: func(mockLoanPaymentRepo *MockLoanPaymentRepo, mockLockManager *MockLockManager) {
				// Mock lock
				mockLockManager.On("GetLock", "loan-id-1").Return(&sync.Mutex{})

				// Mock unpaid loan payments
				loanPayments := []*model.LoanPayment{
					{
						ID:         uuid.Must(uuid.NewV7()).String(),
						LoanID:     "loan-id-1",
						BorrowerID: "borrower-id-1",
						Amount:     decimal.NewFromInt(110_000),
						DueDate:    pastDue,
						Status:     constant.LoanPaymentStatusUnpaid,
					},
					{
						ID:         uuid.Must(uuid.NewV7()).String(),
						LoanID:     "loan-id-1",
						BorrowerID: "borrower-id-1",
						Amount:     decimal.NewFromInt(110_000),
						DueDate:    futureDue,
						Status:     constant.LoanPaymentStatusUnpaid,
					},
				}
				mockLoanPaymentRepo.On("Find", mock.MatchedBy(func(lp model.LoanPayment) bool {
					return lp.LoanID == "loan-id-1" && lp.Status == constant.LoanPaymentStatusUnpaid
				})).Return(loanPayments, nil)

				// Mock change status to paid
				mockLoanPaymentRepo.On("ChangeStatusToPaid", []string{loanPayments[0].ID}, mock.Anything).Return(nil)
			},
			expectedError: false,
		},
		{
			name:   "Error - Payment Less Than Minimum",
			loanID: "loan-id-2",
			amount: decimal.NewFromInt(50_000),
			mockSetup: func(mockLoanPaymentRepo *MockLoanPaymentRepo, mockLockManager *MockLockManager) {
				// Mock lock
				mockLockManager.On("GetLock", "loan-id-2").Return(&sync.Mutex{})

				// Mock unpaid loan payments
				loanPayments := []*model.LoanPayment{
					{
						ID:         uuid.Must(uuid.NewV7()).String(),
						LoanID:     "loan-id-2",
						BorrowerID: "borrower-id-1",
						Amount:     decimal.NewFromInt(110_000),
						DueDate:    pastDue,
						Status:     constant.LoanPaymentStatusUnpaid,
					},
					{
						ID:         uuid.Must(uuid.NewV7()).String(),
						LoanID:     "loan-id-2",
						BorrowerID: "borrower-id-1",
						Amount:     decimal.NewFromInt(110_000),
						DueDate:    futureDue,
						Status:     constant.LoanPaymentStatusUnpaid,
					},
				}
				mockLoanPaymentRepo.On("Find", mock.MatchedBy(func(lp model.LoanPayment) bool {
					return lp.LoanID == "loan-id-2" && lp.Status == constant.LoanPaymentStatusUnpaid
				})).Return(loanPayments, nil)
			},
			expectedError: true,
		},
		{
			name:   "Error - Payment Not In Plan",
			loanID: "loan-id-3",
			amount: decimal.NewFromInt(150_000),
			mockSetup: func(mockLoanPaymentRepo *MockLoanPaymentRepo, mockLockManager *MockLockManager) {
				// Mock lock
				mockLockManager.On("GetLock", "loan-id-3").Return(&sync.Mutex{})

				// Mock unpaid loan payments
				loanPayments := []*model.LoanPayment{
					{
						ID:         uuid.Must(uuid.NewV7()).String(),
						LoanID:     "loan-id-3",
						BorrowerID: "borrower-id-1",
						Amount:     decimal.NewFromInt(110_000),
						DueDate:    pastDue,
						Status:     constant.LoanPaymentStatusUnpaid,
					},
					{
						ID:         uuid.Must(uuid.NewV7()).String(),
						LoanID:     "loan-id-3",
						BorrowerID: "borrower-id-1",
						Amount:     decimal.NewFromInt(110_000),
						DueDate:    futureDue,
						Status:     constant.LoanPaymentStatusUnpaid,
					},
				}
				mockLoanPaymentRepo.On("Find", mock.MatchedBy(func(lp model.LoanPayment) bool {
					return lp.LoanID == "loan-id-3" && lp.Status == constant.LoanPaymentStatusUnpaid
				})).Return(loanPayments, nil)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLoanRepo := new(MockLoanRepo)
			mockLoanPaymentRepo := new(MockLoanPaymentRepo)
			mockLockManager := new(MockLockManager)
			tt.mockSetup(mockLoanPaymentRepo, mockLockManager)

			// We need to use a type assertion here because LoanService expects lib.LockManager
			service := &LoanService{
				loanRepo:        mockLoanRepo,
				loanPaymentRepo: mockLoanPaymentRepo,
				lockManager:     mockLockManager,
			}
			err := service.MakePayment(tt.loanID, tt.amount)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockLoanPaymentRepo.AssertExpectations(t)
			mockLockManager.AssertExpectations(t)
		})
	}
}
