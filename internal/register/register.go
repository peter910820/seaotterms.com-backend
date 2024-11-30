package register

import (
	"errors"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/crud"
	"seaotterms.com-backend/internal/model"
)

type apiAccount struct {
	Username string
	Email    string
}

func Register(formData *map[string]interface{}) error {
	var find []apiAccount
	dsn := crud.InitDsn()
	data := model.Account{Username: (*formData)["username"].(string),
		Password: (*formData)["password"].(string),
		Email:    (*formData)["email"].(string)}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("database access error: %v", err)
	}

	result := db.Model(&model.Account{}).Find(&find)
	if result.Error != nil {
		log.Fatalf("%v\n", result.Error)
	}
	// check Username & Email exist
	for _, col := range find {
		if data.Username == col.Username {
			return errors.New("username is exist")
		} else if data.Email == col.Email {
			return errors.New("email is exist")
		} else {
		}
	}
	result = db.Create(&data)
	if result.Error != nil {
		log.Fatalf("%v\n", result.Error)
	}
	return nil
}
