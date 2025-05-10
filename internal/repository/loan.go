package repository

import (
	"github/ramabmtr/billing-engine/internal/constant"
	"github/ramabmtr/billing-engine/internal/model"
	"gorm.io/gorm"
)

type LoanRepo interface {
	WithTx(tx *gorm.DB) LoanRepo
	Create(l *model.Loan) error
	Get(l *model.Loan) error
	FindByBorrowerID(borrowerID string) ([]*model.LoanWithCompleteStatus, error)
}

type loanRepo struct {
	db *gorm.DB
}

func NewLoanRepo(db *gorm.DB) LoanRepo {
	return &loanRepo{db: db}
}

func (r *loanRepo) WithTx(tx *gorm.DB) LoanRepo {
	return &loanRepo{db: tx}
}

func (r *loanRepo) Create(l *model.Loan) error {
	return r.db.Create(l).Error
}

func (r *loanRepo) Get(l *model.Loan) error {
	return r.db.First(l).Error
}

func (r *loanRepo) FindByBorrowerID(borrowerID string) ([]*model.LoanWithCompleteStatus, error) {
	var loans = make([]*model.LoanWithCompleteStatus, 0)
	err := r.db.
		Select(
			"l.*",
			`case
						when count(lp.id) > 0
							then false
						else true
					end as is_completed`,
		).
		Table("loans l").
		Joins("left join loan_payments lp on lp.loan_id = l.id and lp.status = ?", constant.LoanPaymentStatusUnpaid).
		Where("l.borrower_id = ?", borrowerID).
		Group("l.id").
		Scan(&loans).Error
	return loans, err
}
