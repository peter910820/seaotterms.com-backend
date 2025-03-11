package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/model"
)

type apiAccount struct {
	Username string
	Email    string
}

type RegisterDataForClient struct {
	Username      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	CheckPassword string `json:"checkPassword"`
}

func RegisterHandler(c *fiber.Ctx, db *gorm.DB) error {
	var data RegisterDataForClient
	var find []apiAccount

	if err := c.BodyParser(&data); err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	result := db.Model(&model.User{}).Find(&find)
	if result.Error != nil {
		logrus.Fatalf("%v", result.Error)
	}

	data.Username = strings.ToLower(data.Username)
	data.Email = strings.ToLower(data.Email)
	// check Username & Email exist
	for _, col := range find {
		if data.Username == col.Username {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"msg": "username is exist",
			})
		} else if data.Email == col.Email {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"msg": "email is exist",
			})
		}
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err,
		})
	}
	data.Password = string(hashedPassword)
	dataCreate := model.User{
		Username:   data.Username,
		Password:   data.Password,
		Email:      data.Email,
		CreateName: data.Username,
	}
	result = db.Create(&dataCreate)
	if result.Error != nil {
		logrus.Error(result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": result.Error,
		})
	}

	return c.JSON(fiber.Map{
		"msg": "註冊成功",
	})
}
