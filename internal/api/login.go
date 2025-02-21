package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/model"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx, store *session.Store, data *LoginData, db *gorm.DB) error {
	var databaseData []LoginData

	r := db.Model(&model.Account{}).Find(&databaseData)
	if r.Error != nil {
		logrus.Fatalf("%v\n", r.Error)
	}

	for _, col := range databaseData {
		if data.Username == col.Username {
			logrus.Infof("Username %s try to login", data.Username)
			if data.Password == col.Password {
				// set session
				sess, err := store.Get(c)
				if err != nil {
					logrus.Fatal(err)
				}
				sess.Set("username", data.Username)
				if err := sess.Save(); err != nil {
					logrus.Fatal(err)
				}
				logrus.Infof("Username %s login success", data.Username)
				return nil
			} else {
				logrus.Infof("login error, password not correct")
				return errors.New("login error, password not correct")
			}
		}
	}
	logrus.Infof("user not found")
	return errors.New("user not found")
}
