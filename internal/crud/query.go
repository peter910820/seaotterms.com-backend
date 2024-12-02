package crud

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"seaotterms.com-backend/internal/model"
)

func Query() {
	var accounts model.Account
	dsn := InitDsn()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("連接資料庫失敗: %v", err)
	}
	result := db.Find(&accounts)
	if result.Error != nil {
		logrus.Fatalf("%v", result.Error)
	}
	fmt.Printf("%v\n", result)
}
