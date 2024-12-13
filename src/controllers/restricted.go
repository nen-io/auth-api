package controllers

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

func Restricted(c fiber.Ctx) error {

	slog.Info("Restricted")

	return c.SendString("Restricted\n")
}
