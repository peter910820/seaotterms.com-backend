package api

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/model"
)

type UserDataForClient struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Exp        int       `json:"exp"`
	Management bool      `json:"management"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdateName string    `json:"update_name"`
	Avatar     string    `json:"avatar"`
}

type UserDataForUpdate struct {
	UpdatedAt  time.Time
	UpdateName string
	Avatar     string
}

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

func UpdateUser(c *fiber.Ctx, db *gorm.DB) error {
	// load client data
	var clientData UserDataForClient
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Errorf("%s\n", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	// URL decoding
	id, err := url.QueryUnescape(c.Params("id"))
	if err != nil {
		logrus.Errorf("%s\n", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	// check if form id equal route id
	u, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	if u != uint64(clientData.ID) {
		logrus.Error("ID比對失敗")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "ID比對失敗",
		})
	}

	r := db.Model(&model.User{}).Where("id = ?", id).
		Select("updated_at", "update_name", "avatar").
		Updates(UserDataForUpdate{
			UpdatedAt:  time.Now(),
			UpdateName: clientData.Username,
			Avatar:     clientData.Avatar,
		})
	if r.Error != nil {
		// if record not exist
		if r.Error == gorm.ErrRecordNotFound {
			logrus.Errorf("%s\n", r.Error.Error())
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"msg": r.Error.Error(),
			})
		} else {
			logrus.Fatal(r.Error.Error())
		}
	}
	logrus.Infof("個人資料 %s 更新成功", clientData.Username)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": fmt.Sprintf("資料 %s 更新成功", clientData.Username),
	})
}
