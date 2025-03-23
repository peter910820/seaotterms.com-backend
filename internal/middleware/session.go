package middleware

import (
	"errors"
	"strings"
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

// check user identity
func SessionHandler(store *session.Store, db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		confirmRoutes := map[string]string{
			"/api/verify":         "POST", // handle front-end router verify
			"/api/create-article": "POST",
			"/api/galgame":        "POST",
			"/api/galgame-brand":  "POST",
		}
		confirmRoutesPrefix := map[string]string{
			"/api/galgame/":       "PATCH",
			"/api/galgame-brand/": "PATCH",
			"/api/users/":         "PATCH",
		}
		if isPathIn(c.Path(), c.Method(), confirmRoutes) {
			err := checkLogin(c, store, db)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"msg":      "visitors is not logged in",
					"userData": UserData{},
				})
			}
			return c.Next()
		}
		if isPathPrefix(c.Path(), c.Method(), confirmRoutesPrefix) {
			err := checkLogin(c, store, db)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"msg":      "visitors is not logged in",
					"userData": UserData{},
				})
			}
			return c.Next()
		}
		return c.Next()
	}
}

// 目前沒用到
func AuthenticationManagementHandler(store *session.Store, db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		confirmRoutes := map[string]string{
			"/api/authentication": "POST",
		}

		if isPathIn(c.Path(), c.Method(), confirmRoutes) {
			sess, err := store.Get(c)
			if err != nil {
				logrus.Fatal(err)
			}
			username := sess.Get("username")
			if username == nil {
				logrus.Error("visitors is not logged in")
				return errors.New("visitors is not logged in")
			}

			var userData model.User

			r := db.Where("username = ?", username).First(&userData)
			if r.Error != nil {
				logrus.Fatal(r.Error.Error())
			}

			if !userData.Management {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"msg": "你沒有權限造訪此頁面",
				})
			}
			return c.Next()
		}
		return c.Next()
	}
}

/* utils */

func isPathIn(path string, method string, confirmRoutes map[string]string) bool {
	for key, value := range confirmRoutes {
		if key == path && value == method {
			return true
		}
	}
	return false
}

func isPathPrefix(path string, method string, confirmRoutesPrefix map[string]string) bool {
	for key, value := range confirmRoutesPrefix {
		if strings.HasPrefix(path, key) && value == method {
			return true
		}
	}
	return false
}

func checkLogin(c *fiber.Ctx, store *session.Store, db *gorm.DB) error {
	sess, err := store.Get(c)
	if err != nil {
		logrus.Fatal(err)
	}
	username := sess.Get("username")
	if username == nil {
		logrus.Error("visitors is not logged in")
		return errors.New("visitors is not logged in")
	}

	var userData model.User

	r := db.Where("username = ?", username).First(&userData)
	if r.Error != nil {
		logrus.Fatal(r.Error.Error())
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
