package crud

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"seaotterms.com-backend/model"
)

type apiAccount struct {
	Username string
	Email    string
}

func Register() {
	var data []apiAccount
	dsn := initDsn()

	// data := model.Account{Username: "testUser", Password: "123456", Email: "example@gmail.com"}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("連接資料庫失敗: %v", err)
	}

	result := db.Model(&model.Account{}).Find(&data)
	if result.Error != nil {
		log.Fatalf("%v\n", result.Error)
	}
	for _, col := range data {
		fmt.Printf("%v\n", col)
	}
	// result := db.Create(&data)
	// if result.Error != nil {
	// 	log.Fatalf("%v", result.Error)
	// }
}
