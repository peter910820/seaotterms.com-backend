package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/dto"
	"seaotterms.com-backend/internal/model"
)

func CreateArticle(c *fiber.Ctx, db *gorm.DB) error {
	var clientData dto.ArticleCreateRequest
	if err := c.BodyParser(&clientData); err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	if len(clientData.Tags) > 0 {
		var count int64
		db.Model(&model.Tag{}).Where("name IN ?", clientData.Tags).Count(&count)
		if count != int64(len(clientData.Tags)) {
			logrus.Error("缺少tags，請先建立tags")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"msg": "缺少tags，請先建立tags",
			})
		}
	}

	data := model.Article{
		Title:   clientData.Title,
		Content: clientData.Content,
		Tags:    []model.Tag{},
	}
	for _, tag := range clientData.Tags {
		if !(strings.TrimSpace(tag) == "") {
			data.Tags = append(data.Tags, model.Tag{Name: tag})
		}
	}

	if err := db.Create(&data).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": "建立資料成功",
	})
}

func QueryArticle(c *fiber.Ctx, db *gorm.DB) error {
	var articleData []model.Article
	var err error

	id := c.Query("id")
	if id != "" {
		err = db.First(&articleData, c.Params("articleID")).Error
	} else {
		err = db.Order("created_at desc").Find(&articleData).Error
	}
	if err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	logrus.Info("query articles success")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": articleData,
	})
}
