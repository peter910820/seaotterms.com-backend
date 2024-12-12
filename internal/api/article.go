package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"seaotterms.com-backend/internal/crud"
	"seaotterms.com-backend/internal/model"
)

type ArticleData struct {
	Title    string         `json:"title"`
	Username string         `json:"username"`
	Tags     pq.StringArray `json:"tags" gorm:"type:text[]"`
	Content  string         `json:"content"`
}

func ArticleHandler(c *fiber.Ctx) error {
	var data ArticleData

	if err := c.BodyParser(&data); err != nil {
		logrus.Fatalf("%v", err)
	}
	err := CreateArticle(&data)
	if err != nil {
		// 500
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func CreateArticle(data *ArticleData) error {
	dsn := crud.InitDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("database access error: %v", err)
	}
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

func GetArticle(c *fiber.Ctx) error {
	var articleData []model.Article

	dsn := crud.InitDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("database access error: %v", err)
	}
	result := db.Order("created_at desc").Find(&articleData)
	if result.Error != nil {
		logrus.Fatalf("%v", result.Error)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": articleData,
	})
}

func GetSingleArticle(c *fiber.Ctx) error {
	var articleData model.Article

	dsn := crud.InitDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("database access error: %v", err)
	}
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
	logrus.Debugf("%v", articleData)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": articleData,
	})
}
