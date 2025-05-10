package main

import (
	"log"

	"github/ramabmtr/billing-engine/config"
	"github/ramabmtr/billing-engine/internal/model"
)

func main() {
	config.InitEnv()
	config.InitDB()

	log.Println("Running database migrations...")

	// Auto migrate the schema
	err := config.GetDB().AutoMigrate(
		&model.Borrower{},
		&model.Loan{},
		&model.LoanPayment{},
	)

	if err != nil {
		log.Printf("Failed to migrate database: %s\n", err.Error())
	}

	log.Println("Database migrations completed successfully")
}
