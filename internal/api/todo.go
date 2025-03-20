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

type todoDataForInsert struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	Owner      string     `gorm:"NOT NULL" json:"owner"`
	Topic      string     `gorm:"NOT NULL" json:"topic"`
	Title      string     `gorm:"NOT NULL" json:"title"`
	Status     uint       `gorm:"NOT NULL" json:"status"`
	Deadline   *time.Time `json:"deadline"`
	CreateName string     `gorm:"NOT NULL" json:"createName"`
	UpdateName string     `json:"updateName"`
}

type todoSwitch struct {
	Status     uint   `gorm:"NOT NULL" json:"status"`
	UpdateName string `json:"updateName"`
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

func SwitchTodo(c *fiber.Ctx, db *gorm.DB) error {
	data := model.Todo{}
	result := db.First(&data, c.Params("id"))
	if result.Error != nil {
		// if record not exist
		if result.Error == gorm.ErrRecordNotFound {
			logrus.Info(result.Error)
			return c.SendStatus(fiber.StatusNotFound)
		} else {
			logrus.Fatal(result.Error)
		}
	}
	status := uint(0)
	if data.Status == 0 {
		status = 3
	} else if data.Status == 3 {
		status = 0
	} else {
		logrus.Fatal("切換API出問題了")
	}
	r := db.Model(&model.Todo{}).Where("id = ?", c.Params("id")).
		Select("status", "update_name").
		Updates(todoSwitch{
			Status:     status,
			UpdateName: "root",
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
	logrus.Infof("Todo資料 %s 切換成功", c.Params("id"))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": fmt.Sprintf("Todo資料 %s 切換成功", c.Params("id")),
	})
}
