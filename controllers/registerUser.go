package controllers

import (
	"encoding/json"
	"log/slog"

	"example.com/models"
	"github.com/gofiber/fiber/v3"
)

func RegisterUser(c fiber.Ctx) error {
	// This function will be used to register a user
	//

	body := c.Request().Body()

	user := models.User{}
	err := json.Unmarshal(body, &user)

	if err != nil || !user.Valid() {
		slog.Error("Failed to unmarshal request body")
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("Invalid request body"))
	}

	slog.Info("REQUEST", "Body", body)

	slog.Info("User", "FirstName", user.FirstName, "LastName", user.LastName, "Email", user.Email, "Password", user.Password)

	return c.SendString("Register User")
}
