package controllers

import (
	"example.com/src/models"
	"github.com/gofiber/fiber/v3"
)

func Verify(c fiber.Ctx) error {

	token := c.Params("token")
	id := c.Params("id")

	// Verify user with token
	u := models.User{
		ID: id,
	}
	uFields := []string{"verified", "verification_token"}

	if err := u.GetFields(uFields); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("User not found"))
	}

	if u.Verified {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("User already verified"))
	}

	if u.VerificationToken != token {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("Invalid verification token"))
	}

	if err := u.Update("verified", true); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("failed to verify user"))
	}

	return c.SendString("Verify User with token: " + token + "\n")
}
