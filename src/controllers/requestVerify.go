package controllers

import (
	"example.com/src/models"
	"example.com/src/services"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func RequestEmailVerificationCode(c fiber.Ctx) error {

	// TODO: use Fiber Bind() to bind the request body to a struct

	request := new(models.RequestEmail)
	if err := c.Bind().JSON(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("Invalid request body"))
	}

	user := models.User{
		Email: request.Email,
	}

	verifyFields := []string{"id", "verified"}

	if err := user.GetFields(verifyFields); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("User not found"))
	}

	if user.Verified {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("User already verified"))
	}

	newVerificationToken := uuid.New().String()
	if err := user.Update("verification_token", newVerificationToken); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.MakeError("Failed to update verification token"))
	}

	// Send email to user
	services.SendVerficationEmail(user.Email, user.VerificationToken, user.ID)

	mapResp := map[string]any{
		"success": true,
		"message": "Verification code sent to email",
	}

	return c.JSON(mapResp)
}
