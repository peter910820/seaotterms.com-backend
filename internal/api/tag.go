package api

import (
	"github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"seaotterms.com-backend/internal/crud"
)

type TagData struct {
	ID    uint
	Title string
}

func GetTags(c *fiber.Ctx) error {
	var tagData []TagData

	dsn := crud.InitDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("database access error: %v", err)
	}
	result := db.Table("articles").Select("id", "title").Where("? = ANY(tags)", c.Params("tagName")).Find(&tagData)
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
