package api

// "errors"
// "log"

// "gorm.io/driver/postgres"
// "gorm.io/gorm"

// "github.com/gofiber/fiber/v2"
// "github.com/gofiber/fiber/v2/middleware/session"

// "seaotterms.com-backend/internal/crud"
// "seaotterms.com-backend/internal/model"

type ArticleData struct {
	Title    string `json:"title"`
	Username string `json:"username"`
	Tags     string `json:"tags"`
	Content  string `json:"content"`
}

func createArticle() {
	// var articleData ArticleData

}
