package controllers

import (
	"log/slog"

	"example.com/src/models"
	"github.com/gofiber/fiber/v3"
)

func RegisterUser(c fiber.Ctx) error {

	user := new(models.User)
	if err := c.Bind().JSON(user); err != nil {
		slog.Error("Failed to unmarshal request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("Invalid request body"))
	}

	if err := user.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError(err.Error()))
	}

	if err := user.Create(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.MakeError("Something went wrong"))

	}

	resp := map[string]string{
		"status":  "VERIFY_EMAIL",
		"message": "User created",
		"id":      user.ID,
	}

	return c.JSON(resp)
}
