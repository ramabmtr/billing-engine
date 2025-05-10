package repository

import (
	"time"

	"github.com/shopspring/decimal"
	"github/ramabmtr/billing-engine/internal/constant"
	"github/ramabmtr/billing-engine/internal/model"
	"gorm.io/gorm"
)

type LoanPaymentRepo interface {
	WithTx(tx *gorm.DB) LoanPaymentRepo
	CreateBulk(lps []*model.LoanPayment) error
	GetTotalOutstandingByLoanID(loanID string) (decimal.Decimal, error)
	GetTotalOutstandingByBorrowerID(borrowerID string) (decimal.Decimal, error)
	Find(lp model.LoanPayment) ([]*model.LoanPayment, error)
	ChangeStatusToPaid(loanIds []string, paidAt time.Time) error
}

type loanPaymentRepo struct {
	db *gorm.DB
}

func NewLoanPaymentRepo(db *gorm.DB) LoanPaymentRepo {
	return &loanPaymentRepo{db: db}
}

func (r *loanPaymentRepo) WithTx(tx *gorm.DB) LoanPaymentRepo {
	return &loanPaymentRepo{db: tx}
}

func (r *loanPaymentRepo) CreateBulk(lps []*model.LoanPayment) error {
	return r.db.Create(lps).Error
}

func (r *loanPaymentRepo) GetTotalOutstandingByLoanID(loanID string) (decimal.Decimal, error) {
	var total decimal.Decimal
	err := r.db.Model(&model.LoanPayment{}).
		Where(&model.LoanPayment{
			LoanID: loanID,
			Status: constant.LoanPaymentStatusUnpaid,
		}).
		Select("coalesce(sum(amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *loanPaymentRepo) GetTotalOutstandingByBorrowerID(borrowerID string) (decimal.Decimal, error) {
	var total decimal.Decimal
	err := r.db.Model(&model.LoanPayment{}).
		Where(&model.LoanPayment{
			BorrowerID: borrowerID,
			Status:     constant.LoanPaymentStatusUnpaid,
		}).
		Select("coalesce(sum(amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *loanPaymentRepo) Find(lp model.LoanPayment) ([]*model.LoanPayment, error) {
	var lps = make([]*model.LoanPayment, 0)
	err := r.db.Where(&lp).Order("due_date asc").Find(&lps).Error
	return lps, err
}

func (r *loanPaymentRepo) ChangeStatusToPaid(loanIds []string, paidAt time.Time) error {
	return r.db.Model(&model.LoanPayment{}).
		Where(&model.LoanPayment{
			Status: constant.LoanPaymentStatusUnpaid,
		}).
		Where("id in ?", loanIds).
		Updates(&model.LoanPayment{
			Status: constant.LoanPaymentStatusPaid,
			PaidAt: &paidAt,
		}).Error
}
