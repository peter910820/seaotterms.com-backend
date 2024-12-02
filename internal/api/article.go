package api

import (
	// "errors"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// "github.com/gofiber/fiber/v2"

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
	// var articleData ArticleData
	dsn := crud.InitDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("database access error: %v", err)
	}

	dataCreate := model.Article{
		Title:    data.Title,
		Username: data.Username,
		Tags:     data.Tags,
		Content:  data.Content,
	}
	fmt.Printf("%v\n", dataCreate)
	result := db.Create(&dataCreate)
	if result.Error != nil {
		log.Fatalf("%v\n", result.Error)
	}
	return nil
}
