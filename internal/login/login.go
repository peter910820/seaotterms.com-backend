package login

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// "github.com/gofiber/fiber/v2"
	// "github.com/gofiber/session/v2"

	"seaotterms.com-backend/internal/crud"
	"seaotterms.com-backend/internal/model"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(data *LoginData) error {
	var databaseData []LoginData

	dsn := crud.InitDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("database access error: %v", err)
	}
	r := db.Model(&model.Account{}).Find(&databaseData)
	if r.Error != nil {
		log.Fatalf("%v\n", r.Error)
	}

	for _, col := range databaseData {
		if data.Username == col.Username {
			fmt.Printf("Username %s try to login", data.Username)
			if data.Password == col.Password {

			} else {
				return errors.New("login error, password not correct")
			}
		} else {
		}
	}
	return nil
}
