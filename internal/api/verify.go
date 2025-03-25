package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func Verify(c *fiber.Ctx, store *session.Store) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg":      "驗證成功",
		"userData": c.Locals("userData"),
	})
}

func AuthenticationManagementHandler(c *fiber.Ctx, store *session.Store) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg":      "驗證成功",
		"userData": c.Locals("userData"),
	})
}
