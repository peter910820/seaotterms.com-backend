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

type RegisterData struct {
	Username      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	CheckPassword string `json:"checkPassword"`
}

func Register(data *RegisterData) error {
	var find []apiAccount
	dsn := crud.InitDsn()

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

	dataCreate := model.Account{
		Username: data.Username,
		Password: data.Password,
		Email:    data.Email,
	}
	result = db.Create(&dataCreate)
	if result.Error != nil {
		log.Fatalf("%v\n", result.Error)
	}
	return nil
}
