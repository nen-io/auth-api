package controllers

import (
	"encoding/json"
	"log/slog"
	"strings"

	"example.com/src/models"
	"github.com/gofiber/fiber/v3"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c fiber.Ctx) error {

	body := c.Request().Body()

	loginRequest := LoginRequest{}

	d := json.NewDecoder(strings.NewReader(string(body)))
	d.DisallowUnknownFields()

	if err := d.Decode(&loginRequest); err != nil {
		slog.Error("login decode body", "error", err, "body", body)
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("Invalid login request body"))
	}
	existingUser := models.User{}
	if err := existingUser.ValidateEmail(loginRequest.Email); err != nil {
		slog.Error("email validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("Invalid email"))
	}

	existingUser.GetByEmail(loginRequest.Email)

	// check password
	if err := existingUser.ComparePassword(loginRequest.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("Failed to login"))
	}

	return c.SendString("Login")

}
