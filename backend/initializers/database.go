package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatal("Unable to connect to database")
		return nil
	}
	return db
}

func CommitOrRollback(tx *gorm.DB) {
	if r := recover(); r != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
}
