package crud

import (
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
		logrus.Fatalf("database access error: %v", err)
	}
	result := db.Find(&accounts)
	if result.Error != nil {
		logrus.Fatalf("%v", result.Error)
	}
	logrus.Debugf("%v", result)
}
