package api

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
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
	u64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	u := uint(u64)
	if u != clientData.ID {
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
