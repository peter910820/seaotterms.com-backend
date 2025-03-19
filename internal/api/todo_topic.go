package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"seaotterms.com-backend/internal/model"
)

func QueryTodoTopic(c *fiber.Ctx, db *gorm.DB) error {
	data := []model.TodoTopic{}
	r := db.Order("topic_name asc").Find(&data)
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
	logrus.Info("查詢TodoTopic資料成功")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg":  "查詢TodoTopic資料成功",
		"data": data,
	})
}
