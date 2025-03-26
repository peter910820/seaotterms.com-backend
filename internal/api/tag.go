package api

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/model"
)

type TagData struct {
	ID    uint
	Title string
}

func QueryTags(c *fiber.Ctx, db *gorm.DB) error {
	var tagData []model.Tag

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
	logrus.Info("查詢全部tag資料成功")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": tagData,
	})
}

func QueryTag(c *fiber.Ctx, db *gorm.DB) error {
	var tagData []TagData

	decodeTag, err := url.QueryUnescape(c.Params("tagName"))
	if err != nil {
		logrus.Fatalf("Failed to decode URL: %v", err)
		// 500
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	result := db.Table("articles").Select("id", "title").Where("? = ANY(tags)", decodeTag).Find(&tagData)
	if result.Error != nil {
		// if record not exist
		if result.Error == gorm.ErrRecordNotFound {
			logrus.Error(result.Error)
			return c.SendStatus(fiber.StatusNotFound)
		} else {
			logrus.Error(result.Error)
			// 500
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"msg": result.Error.Error(),
			})
		}
	}
	logrus.Infof("查詢指定標籤 %s 全部文章資料成功", decodeTag)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": tagData,
	})
}
