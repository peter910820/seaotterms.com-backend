package api

import (
	"fmt"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/model"
)

func QueryTodoByOwner(c *fiber.Ctx, db *gorm.DB) error {
	// URL decoding
	owner, err := url.QueryUnescape(c.Params("owner"))
	if err != nil {
		logrus.Error(err)
		// return 400
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	data := []model.Todo{}
	r := db.Where("owner = ?", owner).Find(&data)
	if r.Error != nil {
		// if record not exist
		if r.Error == gorm.ErrRecordNotFound {
			logrus.Error(r.Error)
			// return 404
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"msg": r.Error.Error(),
			})
		} else {
			logrus.Fatal(r.Error.Error())
		}
	}
	logrus.Infof("查詢%s的Todo資料成功", owner)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg":  fmt.Sprintf("查詢 %s 的Todo資料成功", owner),
		"data": data,
	})
}

func InsertTodo(c *fiber.Ctx, db *gorm.DB) error {
	// load client data
	clientData := model.Todo{}
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	r := db.Create(&clientData)
	if r.Error != nil {
		logrus.Errorf("%s\n", r.Error.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": r.Error.Error(),
		})
	}
	logrus.Infof("資料 %s 創建成功", clientData.Title)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": fmt.Sprintf("資料 %s 創建成功", clientData.Title),
	})
}
