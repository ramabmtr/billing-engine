package repository

import (
	"github.com/ramabmtr/billing-engine/internal/constant"
	"github.com/ramabmtr/billing-engine/internal/model"
	"gorm.io/gorm"
)

type BorrowerRepo interface {
	WithTx(tx *gorm.DB) BorrowerRepo
	Create(b *model.Borrower) error
	List() ([]*model.BorrowerWithDelinquentStatus, error)
}

type borrowerRepo struct {
	db *gorm.DB
}

func NewBorrowerRepo(db *gorm.DB) BorrowerRepo {
	return &borrowerRepo{db: db}
}

func (r *borrowerRepo) WithTx(tx *gorm.DB) BorrowerRepo {
	return &borrowerRepo{db: tx}
}

func (r *borrowerRepo) Create(b *model.Borrower) error {
	return r.db.Create(b).Error
}

func (r *borrowerRepo) List() ([]*model.BorrowerWithDelinquentStatus, error) {
	var borrowers = make([]*model.BorrowerWithDelinquentStatus, 0)
	err := r.db.
		Select(
			"b.*",
			`case
	when count(lp.id) > 1
		then true
		else false
		end as is_delinquent`,
		).
		Table("borrowers b").
		Joins("left join loan_payments lp ON lp.borrower_id = b.id and lp.status = ? and lp.due_date < ?", constant.LoanPaymentStatusUnpaid, r.db.NowFunc()).
		Group("b.id").
		Scan(&borrowers).Error
	return borrowers, err
}
