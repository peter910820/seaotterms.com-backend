package model

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func Migration(dbName string, db *gorm.DB) {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf(".env file error: %v", err)
	}

	switch dbName {
	case os.Getenv("DB_NAME"):
		db.AutoMigrate(&Account{})
		db.AutoMigrate(&Article{})
		db.AutoMigrate(&Tag{})
	case os.Getenv("DB_NAME2"):
		db.AutoMigrate(&BrandRecord{})
		db.AutoMigrate(&GameRecord{})
	case os.Getenv("DB_NAME3"):
		db.AutoMigrate(&User{})
	default:
		logrus.Fatal("error in migration function")
	}
}
