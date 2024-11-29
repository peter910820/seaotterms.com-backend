package crud

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "seaotterms.com-backend/model"
)

func Register() {
	// data := model.Account{Username: "testUser", Password: "123456", Email: "example@gmail.com"}

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_OWNER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))
	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("連接資料庫失敗: %v", err)
	}

	// result := db.Create(&data)
	// if result.Error != nil {
	// 	log.Fatalf("%v", result.Error)
	// }
}
