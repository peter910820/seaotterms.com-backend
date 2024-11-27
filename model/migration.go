package model

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Migration() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(".env file error: %v", err)
	}
	// PostgreSQL 連線字串
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_OWNER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("連接資料庫失敗: %v", err)
	}
	db.AutoMigrate(&Account{})
}
