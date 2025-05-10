package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ramabmtr/billing-engine/config"
	"github.com/ramabmtr/billing-engine/internal/constant"
	"github.com/ramabmtr/billing-engine/internal/lib"
	"github.com/ramabmtr/billing-engine/internal/model"
	"github.com/ramabmtr/billing-engine/internal/repository"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type LoanService struct {
	loanRepo        repository.LoanRepo
	loanPaymentRepo repository.LoanPaymentRepo
	lockManager     lib.LockManager
}

func NewLoanService(loanRepo repository.LoanRepo, loanPaymentRepo repository.LoanPaymentRepo) *LoanService {
	return &LoanService{
		loanRepo:        loanRepo,
		loanPaymentRepo: loanPaymentRepo,
		lockManager:     lib.NewLockManager(),
	}
}

func (s *LoanService) CreateLoanRequest(borrowerID string) (*model.Loan, error) {
	// check if there is an outstanding amount for that borrower id
	outstandingAmount, err := s.loanPaymentRepo.GetTotalOutstandingByBorrowerID(borrowerID)
	if err != nil {
		return nil, err
	}
	if !outstandingAmount.IsZero() {
		return nil, fmt.Errorf("there is an outstanding loan for this borrower")
	}

	l := &model.Loan{
		ID:                 uuid.Must(uuid.NewV7()).String(),
		BorrowerID:         borrowerID,
		Principal:          decimal.NewFromInt(5_000_000),
		AnnualInterestRate: decimal.NewFromInt(10),
		Period:             50,
		PeriodUnit:         constant.PeriodUnitWeek,
		CreatedAt:          time.Now().UTC(),
	}

	l.TotalRepayment = lib.CalculateTotalRepayment(l.Principal, l.AnnualInterestRate, l.Period, l.PeriodUnit).Round(0)

	err = config.GetDB().Transaction(func(tx *gorm.DB) error {
		err := s.loanRepo.WithTx(tx).Create(l)
		if err != nil {
			return err
		}
		err = s.loanPaymentRepo.WithTx(tx).CreateBulk(s.generateLoanPayment(*l))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (s *LoanService) generateLoanPayment(l model.Loan) []*model.LoanPayment {
	var lps = make([]*model.LoanPayment, l.Period)
	for i := 0; i < l.Period; i++ {
		lps[i] = &model.LoanPayment{
			LoanID:     l.ID,
			BorrowerID: l.BorrowerID,
			Amount:     l.TotalRepayment.Div(decimal.NewFromInt(int64(l.Period))),
			DueDate:    l.CreatedAt.AddDate(0, 0, 7*(i+1)),
			Status:     constant.LoanPaymentStatusUnpaid,
		}
	}

	return lps
}

func (s *LoanService) GetLoansByBorrowerID(borrowerID string) ([]*model.LoanWithCompleteStatus, error) {
	ls, err := s.loanRepo.FindByBorrowerID(borrowerID)
	if err != nil {
		return nil, err
	}

	return ls, nil
}

func (s *LoanService) GetLoanDetail(id string) (*model.Loan, decimal.Decimal, error) {
	l := &model.Loan{
		ID: id,
	}
	err := s.loanRepo.Get(l)
	if err != nil {
		return nil, decimal.NewFromInt(0), err
	}

	o, err := s.loanPaymentRepo.GetTotalOutstandingByLoanID(id)
	if err != nil {
		return nil, decimal.NewFromInt(0), err
	}

	return l, o, nil
}

func (s *LoanService) GetLoanPaymentsByLoanID(loanID string) ([]*model.LoanPayment, error) {
	lps, err := s.loanPaymentRepo.Find(model.LoanPayment{
		LoanID: loanID,
	})
	if err != nil {
		return nil, err
	}

	return lps, nil
}

func (s *LoanService) MakePayment(loanID string, amount decimal.Decimal) error {
	lock := s.lockManager.GetLock(loanID)
	lock.Lock()
	defer lock.Unlock()

	lps, err := s.loanPaymentRepo.Find(model.LoanPayment{
		LoanID: loanID,
		Status: constant.LoanPaymentStatusUnpaid,
	})
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	minimumPayment := decimal.NewFromInt(0)
	paymentPlan := make([]decimal.Decimal, 0)
	tempAmount := decimal.NewFromInt(0)
	for _, lp := range lps {
		tempAmount = tempAmount.Add(lp.Amount)
		paymentPlan = append(paymentPlan, tempAmount)
		if lp.DueDate.Before(now) {
			minimumPayment = minimumPayment.Add(lp.Amount)
		}
	}

	isInPlan := false
	planIndex := -1
	for i, plan := range paymentPlan {
		if amount.Equal(plan) {
			isInPlan = true
			planIndex = i
			break
		}
	}

	if amount.LessThan(minimumPayment) {
		return fmt.Errorf("you must make payment equal to %s at minimum", minimumPayment)
	}

	if !isInPlan {
		return fmt.Errorf("you must make payment equal to %s at minimum or multiples thereof and maximum %s", paymentPlan[0], paymentPlan[len(paymentPlan)-1])
	}

	idToUpdate := make([]string, planIndex+1)
	for i := 0; i <= planIndex; i++ {
		idToUpdate[i] = lps[i].ID
	}
	err = s.loanPaymentRepo.ChangeStatusToPaid(idToUpdate, now)
	if err != nil {
		return err
	}

	return nil
}
