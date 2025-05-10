package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github/ramabmtr/billing-engine/internal/constant"
	"github/ramabmtr/billing-engine/internal/lib"
	"gorm.io/gorm"
)

type Loan struct {
	ID                 string                  `json:"id" gorm:"type:char(36);primary_key"`
	BorrowerID         string                  `json:"borrower_id" gorm:"type:char(36);not null"`
	Borrower           *Borrower               `json:"borrower,omitempty" gorm:"foreignKey:BorrowerID;references:ID"`
	Principal          decimal.Decimal         `json:"principal" gorm:"type:decimal(16,4);not null"`
	AnnualInterestRate decimal.Decimal         `json:"annual_interest_rate" gorm:"type:decimal(5,2);not null"`
	TotalRepayment     decimal.Decimal         `json:"total_repayment" gorm:"type:decimal(16,4);not null"`
	Period             int                     `json:"period" gorm:"type:integer;not null"`
	PeriodUnit         constant.LoanPeriodUnit `json:"period_unit" gorm:"type:varchar(5);not null"`
	CreatedAt          time.Time               `json:"created_at" gorm:"type:timestamp;default:now();not null"`
}

func (c *Loan) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.Must(uuid.NewV7()).String()
	}
	if c.TotalRepayment.IsZero() {
		c.TotalRepayment = lib.CalculateTotalRepayment(c.Principal, c.AnnualInterestRate, c.Period, c.PeriodUnit).Round(0)
	}
	return nil
}

type LoanWithCompleteStatus struct {
	Loan
	IsCompleted bool `json:"is_completed"`
}
