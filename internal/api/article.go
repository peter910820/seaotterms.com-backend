package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/model"
)

type ArticleData struct {
	Title    string         `json:"title"`
	Username string         `json:"username"`
	Tags     pq.StringArray `json:"tags" gorm:"type:text[]"`
	Content  string         `json:"content"`
}

func CreateArticle(c *fiber.Ctx, db *gorm.DB) error {
	var data ArticleData
	if err := c.BodyParser(&data); err != nil {
		logrus.Error(err)
		// 500
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	dataCreate := model.Article{
		Title:    data.Title,
		Username: data.Username,
		Tags:     data.Tags,
		Content:  data.Content,
	}
	logrus.Infof("A article has been create, title name: %s", dataCreate.Title)
	r := db.Create(&dataCreate)
	if r != nil {
		logrus.Error(r.Error)
		// 500
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": r.Error.Error(),
		})
	}
	// check if tag not exist
	var existTags []model.Tag
	r = db.Where("name = ANY(?)", dataCreate.Tags).Find(&existTags)
	if r != nil {
		logrus.Error(r.Error)
		// 500
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": r.Error.Error(),
		})
	}

	existTagsNames := make(map[string]bool)
	for _, tag := range existTags {
		existTagsNames[tag.Name] = true
	}

	var newTags []model.Tag
	for _, tagName := range data.Tags {
		if !existTagsNames[tagName] {
			newTags = append(newTags, model.Tag{Name: tagName})
		}
	}

	if len(newTags) > 0 {
		r := db.Create(&newTags)
		if r.Error != nil {
			logrus.Error(r.Error)
			// 500
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"msg": r.Error.Error(),
			})
		}
	}
	return nil
}

func QueryArticle(c *fiber.Ctx, db *gorm.DB) error {
	var articleData []model.Article

	result := db.Order("created_at desc").Find(&articleData)
	if result.Error != nil {
		logrus.Error(result.Error)
		// 500
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": result.Error.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": articleData,
	})
}

func QuerySingleArticle(c *fiber.Ctx, db *gorm.DB) error {
	var articleData model.Article

	// find articles
	result := db.First(&articleData, c.Params("articleID"))
	if result.Error != nil {
		// if record not exist
		if result.Error == gorm.ErrRecordNotFound {
			logrus.Info(result.Error)
			return c.SendStatus(fiber.StatusNotFound)
		} else {
			// 500
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"msg": result.Error.Error(),
			})
		}
	}
	logrus.Info("查詢單一文章成功")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": articleData,
	})
}
