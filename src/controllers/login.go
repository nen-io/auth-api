package controllers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"strings"

	"example.com/src/models"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
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

	// only fetch what we need for login
	user := models.User{
		Email: loginRequest.Email,
	}
	userLoginFields := []string{"password", "verified"}

	if err := user.ValidateEmail(loginRequest.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("Invalid email"))
	}

	if err := user.GetFields(userLoginFields); errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("User not found"))
	}

	// check password
	if err := user.ComparePassword(loginRequest.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("Failed to login"))
	}

	if user.Verified == false {
		resp := map[string]string{
			"status":  "VERIFY EMAIL",
			"message": "User not verified",
		}

		return c.JSON(resp)
	}

	return c.SendString("Login")

}
