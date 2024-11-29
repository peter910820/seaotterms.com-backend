package crud

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"seaotterms.com-backend/model"
)

func Query() {
	var accounts model.Account
	dsn := initDsn()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("連接資料庫失敗: %v", err)
	}
	result := db.Find(&accounts)
	if result.Error != nil {
		log.Fatalf("%v", result.Error)
	}
	fmt.Printf("%v\n", result)
}
