package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Borrower struct {
	ID        string    `json:"id" gorm:"type:char(36);primary_key"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;default:now();not null"`
}

func (c *Borrower) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.Must(uuid.NewV7()).String()
	}
	return nil
}

type BorrowerWithDelinquentStatus struct {
	Borrower
	IsDelinquent bool `json:"is_delinquent"`
}
