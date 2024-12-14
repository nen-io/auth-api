package controllers

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

func Restricted(c fiber.Ctx) error {

	slog.Info("Restricted", "id", c.Locals("userId"), "email", c.Locals("email"))

	return c.SendString("Restricted\n")
}
