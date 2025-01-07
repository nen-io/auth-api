package controllers

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

func Me(c fiber.Ctx) error {

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
