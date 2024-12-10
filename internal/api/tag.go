package api

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/crud"
	"seaotterms.com-backend/internal/model"
)

type TagData struct {
	ID    uint
	Title string
}

func GetTags(c *fiber.Ctx) error {
	var tagData []model.Tag

	dsn := crud.InitDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("database access error: %v", err)
	}
	result := db.Order("id desc").Find(&tagData)
	if result.Error != nil {
		// if record not exist
		if result.Error == gorm.ErrRecordNotFound {
			logrus.Info(result.Error)
			return c.SendStatus(fiber.StatusNotFound)
		} else {
			logrus.Fatal(result.Error)
		}
	}
	logrus.Debugf("%v", tagData)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": tagData,
	})
}

func GetTag(c *fiber.Ctx) error {
	var tagData []TagData

	dsn := crud.InitDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("database access error: %v", err)
	}
	decodeTag, err := url.QueryUnescape(c.Params("tagName"))
	if err != nil {
		logrus.Fatalf("Failed to decode URL: %v", err)
	}
	result := db.Table("articles").Select("id", "title").Where("? = ANY(tags)", decodeTag).Find(&tagData)
	if result.Error != nil {
		// if record not exist
		if result.Error == gorm.ErrRecordNotFound {
			logrus.Info(result.Error)
			return c.SendStatus(fiber.StatusNotFound)
		} else {
			logrus.Fatal(result.Error)
		}
	}
	logrus.Debugf("%v", tagData)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": tagData,
	})
}
