package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"seaotterms.com-backend/internal/model"
)

type todoTopicForInsert struct {
	TopicName  string `gorm:"primaryKey" json:"topicName"`
	TopicOwner string `gorm:"NOT NULL; default:'common'" json:"topicOwner"`
	UpdateName string `json:"updateName"`
}

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

func InsertTodoTopic(c *fiber.Ctx, db *gorm.DB) error {
	// load client data
	clientData := todoTopicForInsert{}
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	data := model.TodoTopic{
		TopicName:  clientData.TopicName,
		TopicOwner: clientData.TopicOwner,
		UpdateName: clientData.UpdateName,
	}
	r := db.Create(&data)
	if r.Error != nil {
		logrus.Errorf("%s\n", r.Error.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": r.Error.Error(),
		})
	}
	logrus.Infof("資料 %s 創建成功", clientData.TopicName)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": fmt.Sprintf("資料 %s 創建成功", clientData.TopicName),
	})
}
