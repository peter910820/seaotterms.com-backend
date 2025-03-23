package api

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/model"
)

type todoForClient struct {
	Status     uint      `json:"status"`
	UpdatedAt  time.Time `json:"updatedAt"`
	UpdateName string    `json:"updateName"`
}

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
	r := db.Where("owner = ?", owner).Order("created_at DESC").Find(&data)
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
	data := model.Todo{
		ID:         clientData.ID,
		Owner:      clientData.Owner,
		Topic:      clientData.Topic,
		Title:      clientData.Title,
		Status:     clientData.Status,
		Deadline:   clientData.Deadline,
		CreateName: clientData.CreateName,
		UpdateName: clientData.UpdateName,
	}
	r := db.Create(&data)
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

func UpdateTodoStatus(c *fiber.Ctx, db *gorm.DB) error {
	// load client data
	var clientData todoForClient
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	clientData.UpdatedAt = time.Now()
	r := db.Model(&model.Todo{}).Where("id = ?", c.Params("id")).
		Select("status", "updated_at", "update_name").
		Updates(clientData)
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
	logrus.Infof("Todo %s 更新成功", c.Params("id"))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": fmt.Sprintf("Todo %s 更新成功", c.Params("id")),
	})
}

func DeleteTodo(c *fiber.Ctx, db *gorm.DB) error {
	r := db.Where("id = ?", c.Params("id")).Delete(&model.Todo{})
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
	logrus.Infof("Todo %s 刪除成功", c.Params("id"))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": fmt.Sprintf("Todo %s 刪除成功", c.Params("id")),
	})
}
