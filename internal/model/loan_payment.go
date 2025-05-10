package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/ramabmtr/billing-engine/internal/constant"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type LoanPayment struct {
	ID         string                     `json:"id" gorm:"type:char(36);primary_key"`
	LoanID     string                     `json:"loan_id" gorm:"type:char(36);not null"`
	Loan       *Loan                      `json:"loan,omitempty" gorm:"foreignKey:LoanID;references:ID"`
	BorrowerID string                     `json:"borrower_id" gorm:"type:char(36);not null"`
	Borrower   *Borrower                  `json:"borrower,omitempty" gorm:"foreignKey:BorrowerID;references:ID"`
	Amount     decimal.Decimal            `json:"amount" gorm:"type:decimal(16,4);not null"`
	DueDate    time.Time                  `json:"due_date" gorm:"type:timestamp;not null"`
	Status     constant.LoanPaymentStatus `json:"status" gorm:"type:varchar(10);not null"`
	PaidAt     *time.Time                 `json:"paid_at" gorm:"type:timestamp;default:null"`
	CreatedAt  time.Time                  `json:"created_at" gorm:"type:timestamp;default:now();not null"`
}

func (c *LoanPayment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.Must(uuid.NewV7()).String()
	}
	return nil
}
