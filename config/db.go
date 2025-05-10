package config

import (
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dbOnce sync.Once
	db     *gorm.DB
)

func InitDB() {
	var err error
	dbOnce.Do(func() {
		dsn := GetEnv().Database.GetDSN()
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			SkipDefaultTransaction: true,
			Logger:                 logger.Default.LogMode(logger.Info),
		})
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %s\n", err.Error())
	}
}

func GetDB() *gorm.DB {
	if db == nil {
		InitDB()
	}
	return db
}
