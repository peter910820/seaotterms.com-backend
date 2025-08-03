package api

import (
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"seaotterms.com-backend/dto"
	"seaotterms.com-backend/model"
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
		err = db.Preload("Tags").First(&articleData, id).Error
	} else {
		err = db.Preload("Tags").Order("created_at desc").Find(&articleData).Error
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

func QueryArticleForTag(c *fiber.Ctx, db *gorm.DB) error {
	var articles []model.Article
	// URL decoding
	name, err := url.QueryUnescape(c.Params("tagName"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	err = db.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
		Joins("JOIN tags ON tags.name = article_tags.tag_name").
		Where("tags.name = ?", name).
		Find(&articles).Error
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": articles,
	})
}
