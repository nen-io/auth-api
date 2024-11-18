package controllers

import (
	"example.com/src/models"
	"github.com/gofiber/fiber/v3"
)

func Verify(c fiber.Ctx) error {

	token := c.Params("token")
	email := c.Params("email")

	// Verify user with token
	u := models.User{
		Email: email,
	}
	uFields := []string{"verified", "verification_token"}

	if err := u.GetFields(uFields); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("User not found"))
	}

	return c.SendString("Verify User with token: " + token + "\n")
}
