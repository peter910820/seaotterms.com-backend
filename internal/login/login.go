package login

import (
	"errors"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	"seaotterms.com-backend/internal/crud"
	"seaotterms.com-backend/internal/model"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx, store *session.Store, data *LoginData) error {
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
			log.Printf("Username %s try to login\n", data.Username)
			if data.Password == col.Password {
				// set session
				sess, err := store.Get(c)
				if err != nil {
					log.Fatal(err.Error())
				}
				sess.Set("username", data.Username)
				if err := sess.Save(); err != nil {
					log.Fatalln(err.Error())
				}
				return nil
			} else {
				return errors.New("login error, password not correct")
			}
		}
	}
	return errors.New("user not found")
}
