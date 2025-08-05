package controllers

import (
	"log/slog"

	"example.com/src/models"
	"github.com/gofiber/fiber/v3"
)

func Me(c fiber.Ctx) error {

	if c.Locals("userId") == nil {
		u := models.User{Email: c.Locals("email").(string)}
		if err := u.GetField("id"); err != nil {
			slog.Error("Failed to get user fields", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.MakeError("Failed to retrieve user information"))
		}
		c.Locals("userId", u.ID)
	}

	slog.Info("Restricted", "id", c.Locals("userId"), "email", c.Locals("email"))

	resp := map[string]any{
		"success":  true,
		"message":  "User verified",
		"email":    c.Locals("email"),
		"id":       c.Locals("userId"),
		"username": c.Locals("username"),
	}

	return c.JSON(resp)
}
