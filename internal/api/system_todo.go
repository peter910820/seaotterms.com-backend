package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/dto"
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

func CreateSystemTodo(c *fiber.Ctx, db *gorm.DB) error {
	// load client data
	var clientData dto.SystemTodoCreateRequest
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	data := model.SystemTodo{
		SystemName:  clientData.SystemName,
		Title:       clientData.Title,
		Detail:      clientData.Detail,
		Status:      clientData.Status,
		Deadline:    clientData.Deadline,
		Urgency:     clientData.Urgency,
		CreatedName: clientData.CreatedName,
	}
	r := db.Create(&data)
	if r.Error != nil {
		logrus.Error(r.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": r.Error.Error(),
		})
	}
	logrus.Infof("系統代辦資料 %s 創建成功", clientData.Title)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": fmt.Sprintf("系統代辦資料 %s 創建成功", clientData.Title),
	})
}

func UpdateSystemTodo(c *fiber.Ctx, db *gorm.DB) error {
	// load client data
	var clientData dto.SystemTodoUpdateRequest
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	updateData := dto.SystemTodoUpdate{
		SystemName:  clientData.SystemName,
		Title:       clientData.Title,
		Detail:      clientData.Detail,
		Status:      clientData.Status,
		Deadline:    clientData.Deadline,
		Urgency:     clientData.Urgency,
		UpdatedName: clientData.UpdatedName,
	}
	// clientData.UpdatedAt = time.Now()
	r := db.Model(&model.SystemTodo{}).Where("id = ?", c.Params("id")).
		Select("system_name", "title", "detail", "status", "deadline", "urgency", "updated_name").
		Updates(updateData)
	if r.Error != nil {
		logrus.Error(r.Error)
		// if record not exist
		if r.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"msg": r.Error.Error(),
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"msg": r.Error.Error(),
			})
		}
	}
	logrus.Infof("SystemTodo %s 更新成功", c.Params("id"))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": fmt.Sprintf("SystemTodo %s 更新成功", c.Params("id")),
	})
}

func DeleteSystemTodo(c *fiber.Ctx, db *gorm.DB) error {
	r := db.Where("id = ?", c.Params("id")).Delete(&model.SystemTodo{})
	if r.Error != nil {
		logrus.Error(r.Error)
		// if record not exist
		if r.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"msg": r.Error.Error(),
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"msg": r.Error.Error(),
			})
		}
	}
	logrus.Infof("SystemTodo %s 刪除成功", c.Params("id"))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": fmt.Sprintf("SystemTodo %s 刪除成功", c.Params("id")),
	})
}
