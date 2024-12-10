package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
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

func RegisterHandler(c *fiber.Ctx) error {
	var data RegisterData

	if err := c.BodyParser(&data); err != nil {
		logrus.Errorf("%v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	logrus.Debugf("Received data: %+v", data)

	err := register(&data)
	if err != nil {
		logrus.Infof("%v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"msg": "註冊成功",
	})
}

func register(data *RegisterData) error {
	var find []apiAccount
	dsn := crud.InitDsn()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("database access error: %v", err)
	}

	result := db.Model(&model.Account{}).Find(&find)
	if result.Error != nil {
		logrus.Fatalf("%v", result.Error)
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
		logrus.Errorf("%v", result.Error)
		return result.Error
	}
	return nil
}
