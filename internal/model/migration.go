package model

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Migration() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf(".env file error: %v", err)
	}

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_OWNER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("連接資料庫失敗: %v", err)
	}
	db.AutoMigrate(&Account{})
	db.AutoMigrate(&Article{})
	db.AutoMigrate(&Tag{})
}
