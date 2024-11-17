package controllers

import (
	"encoding/json"
	"log/slog"
	"strings"

	"example.com/src/models"
	"github.com/gofiber/fiber/v3"
)

func RegisterUser(c fiber.Ctx) error {
	// This function will be used to register a user
	//

	body := c.Request().Body()

	user := models.User{}

	d := json.NewDecoder(strings.NewReader(string(body)))
	d.DisallowUnknownFields()

	if err := d.Decode(&user); err != nil {
		slog.Error("Failed to unmarshal request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("Invalid request body"))
	}

	if err := user.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError(err.Error()))
	}

	db_err := user.Save()

	if db_err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.MakeError("Something went wrong"))
	}

	slog.Info("User", "FirstName", user.FirstName, "LastName", user.LastName, "Email", user.Email, "Password", user.Password)

	return c.SendString("Register User")
}
