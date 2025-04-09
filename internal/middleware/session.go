package middleware

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"seaotterms.com-backend/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type UserData struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Exp        int       `json:"exp"`
	Management bool      `json:"management"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdateName string    `json:"update_name"`
	Avatar     string    `json:"avatar"`
}

func CheckLogin(store *session.Store, db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username, err := checkLogin(c, store)
		if err != nil {
			logrus.Warn("visitors is not logged in")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"msg": "visitors is not logged in",
			})
		}
		strUsername, ok := username.(string)
		if !ok {
			logrus.Fatal("session type is not a string, please check the code")
		}
		err = refreshProfile(c, strUsername, db)
		if err != nil {
			logrus.Error("session has error")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"msg": "session has error",
			})
		}
		return c.Next()
	}
}

// check if user is root or seaotterms
func CheckOwner(store *session.Store, db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username, err := checkLogin(c, store)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"msg": "visitors is not logged in",
			})
		}
		strUsername, ok := username.(string)
		if !ok {
			logrus.Fatal("session type is not a string, please check the code")
		}
		err = refreshProfile(c, strUsername, db)
		if err != nil {
			logrus.Error("session has error")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"msg": "session has error",
			})
		}
		if strUsername != "root" && strUsername != "seaotterms" {
			logrus.Error("你沒有權限造訪此頁面")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"msg":      "你沒有權限造訪此頁面",
				"userData": c.Locals("userData"),
			})
		}
		return c.Next()
	}
}

/* utils */

func checkLogin(c *fiber.Ctx, store *session.Store) (interface{}, error) {
	sess, err := store.Get(c)
	if err != nil {
		logrus.Fatal(err)
	}
	username := sess.Get("username")
	if username == nil {
		return nil, errors.New("visitors is not logged in")
	}
	return username, nil
}

func refreshProfile(c *fiber.Ctx, username string, db *gorm.DB) error {
	var userData model.User

	r := db.Where("username = ?", username).First(&userData)
	if r.Error != nil {
		return r.Error
	}

	data := UserData{
		ID:         userData.ID,
		Username:   userData.Username,
		Email:      userData.Email,
		Exp:        userData.Exp,
		Management: userData.Management,
		CreatedAt:  userData.CreatedAt,
		UpdatedAt:  userData.UpdatedAt,
		UpdateName: userData.UpdateName,
		Avatar:     userData.Avatar,
	}
	c.Locals("userData", data)
	logrus.Infof("%s is access %s", username, c.Path())
	return nil
}
