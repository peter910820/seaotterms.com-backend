package api

import (
	"net/url"

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

func ArticleHandler(c *fiber.Ctx, db *gorm.DB) error {
	var data ArticleData

	if err := c.BodyParser(&data); err != nil {
		logrus.Fatalf("%v", err)
	}
	err := createArticle(&data, db)
	if err != nil {
		// 500
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func createArticle(data *ArticleData, db *gorm.DB) error {
	dataCreate := model.Article{
		Title:    data.Title,
		Username: data.Username,
		Tags:     data.Tags,
		Content:  data.Content,
	}
	logrus.Infof("A article has been create, title name: %s", dataCreate.Title)
	result := db.Create(&dataCreate)
	if result.Error != nil {
		logrus.Fatalf("%v\n", result.Error)
	}

	// check if tag not exist
	var existTags []model.Tag
	db.Where("name = ANY(?)", dataCreate.Tags).Find(&existTags)

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
		result := db.Create(&newTags)
		if result.Error != nil {
			logrus.Fatalf("%v\n", result.Error)
		}
	}
	return nil
}

func GetArticle(c *fiber.Ctx, db *gorm.DB) error {
	var articleData []model.Article

	result := db.Order("created_at desc").Find(&articleData)
	if result.Error != nil {
		logrus.Fatalf("%v", result.Error)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": articleData,
	})
}

func GetSingleArticle(c *fiber.Ctx, db *gorm.DB) error {
	var articleData model.Article

	// find articles
	result := db.First(&articleData, c.Params("articleID"))
	if result.Error != nil {
		// if record not exist
		if result.Error == gorm.ErrRecordNotFound {
			logrus.Info(result.Error)
			return c.SendStatus(fiber.StatusNotFound)
		} else {
			logrus.Fatal(result.Error)
		}
	}
	logrus.Info("查詢單一文章成功")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": articleData,
	})
}

func GetTag(c *fiber.Ctx, db *gorm.DB) error {
	var tagData []TagData

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
	logrus.Infof("查詢指定標籤 %s 全部文章資料成功", decodeTag)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": tagData,
	})
}
