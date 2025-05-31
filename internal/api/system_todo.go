package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/model"
)

func QuerySystemTodo(c *fiber.Ctx, db *gorm.DB) error {
	// get query param
	systemName := c.Query("system_name")

	var data []model.SystemTodo
	var r *gorm.DB
	if systemName == "" {
		r = db.Order("COALESCE(updated_at, created_at) DESC").Find(&data)
	} else {
		r = db.Where("system_name = ?", systemName).Order("COALESCE(updated_at, created_at) DESC").Find(&data)
	}
	if r.Error != nil {
		// if record not exist
		if r.Error == gorm.ErrRecordNotFound {
			logrus.Error(r.Error)
			//404
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"msg": r.Error.Error(),
			})
		} else {
			logrus.Fatal(r.Error.Error())
			// 500
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"msg": r.Error.Error(),
			})
		}
	}

	logrus.Info("Query system_Todos table success")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg":  "查詢SystemTodo資料成功",
		"data": data,
	})
}
