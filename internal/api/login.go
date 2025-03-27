package api

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/middleware"
	"seaotterms.com-backend/internal/model"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx, store *session.Store, db *gorm.DB) error {
	var data LoginData
	var databaseData []LoginData

	if err := c.BodyParser(&data); err != nil {
		logrus.Fatal(err)
	}

	r := db.Model(&model.User{}).Find(&databaseData)
	if r.Error != nil {
		logrus.Fatalf("%v\n", r.Error)
	}

	data.Username = strings.ToLower(data.Username)
	for _, col := range databaseData {
		if data.Username == col.Username {
			logrus.Infof("Username %s try to login", data.Username)
			err := bcrypt.CompareHashAndPassword([]byte(col.Password), []byte(data.Password))
			if err != nil {
				logrus.Error("login error, password not correct")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"msg": "密碼輸入錯誤"})
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

			var userData model.User

			r := db.Where("username = ?", data.Username).First(&userData)
			if r.Error != nil {
				logrus.Fatal(r.Error.Error())
			}

			data := middleware.UserData{
				ID:         userData.ID,
				Username:   userData.Username,
				Email:      userData.Email,
				Exp:        userData.Exp,
				Management: userData.Management,
				CreatedAt:  userData.CreatedAt,
				UpdatedAt:  userData.UpdatedAt,
				UpdateName: userData.UpdateName,
				Avatar:     userData.Avatar,
			}
			return c.JSON(fiber.Map{
				"msg":      fmt.Sprintf("%s 登入成功", data.Username),
				"userData": data,
			})
		}
	}
	logrus.Error("user not found")
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"msg": "找不到該使用者"})
}

/* utils */

func CheckPassword(hashedPassword, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err == nil
}
