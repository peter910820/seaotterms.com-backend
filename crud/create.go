package crud

import (
	// "fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "seaotterms.com-backend/model"
)

func Register() {
	// data := model.Account{Username: "testUser", Password: "123456", Email: "example@gmail.com"}
	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("連接資料庫失敗: %v", err)
	}

	// result := db.Create(&data)
	// if result.Error != nil {
	// 	log.Fatalf("%v", result.Error)
	// }
}
