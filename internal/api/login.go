package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/model"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx, store *session.Store, data *LoginData, db *gorm.DB) error {
	var databaseData []LoginData

	r := db.Model(&model.User{}).Find(&databaseData)
	if r.Error != nil {
		logrus.Fatalf("%v\n", r.Error)
	}

	for _, col := range databaseData {
		if data.Username == col.Username {
			logrus.Infof("Username %s try to login", data.Username)
			err := bcrypt.CompareHashAndPassword([]byte(col.Password), []byte(data.Password))
			if err != nil {
				logrus.Infof("login error, password not correct")
				return errors.New("login error, password not correct")
			}
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
		}
	}
	logrus.Infof("user not found")
	return errors.New("user not found")
}

func CheckPassword(hashedPassword, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err == nil // err 為 nil 代表密碼匹配成功
}
