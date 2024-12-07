package api

import (
	"github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"

	"github.com/lib/pq"
	"seaotterms.com-backend/internal/crud"
	"seaotterms.com-backend/internal/model"
)

type ArticleData struct {
	Title    string         `json:"title"`
	Username string         `json:"username"`
	Tags     pq.StringArray `json:"tags" gorm:"type:text[]"`
	Content  string         `json:"content"`
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
