package controllers

import (
	"example.com/src/models"
	"example.com/src/services"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func ForgotPassword(c fiber.Ctx) error {
	request := new(models.ForgotPassword)

	if err := c.Bind().JSON(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("Invalid request body"))
	}

	user := models.User{
		Email: request.Email,
	}

	exists, err := user.EmailExists(user.Email)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.MakeError("Something went wrong"))
	}

	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("User not found"))
	}

	if err := user.GetField("resetPasswordToken"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("Something went wrong"))
	}

	newVerificationToken := uuid.New().String()
	if err := user.Update("resetPasswordToken", newVerificationToken); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.MakeError("Failed to update reset password token"))
	}

	// Send email to user
	if err := services.SendResetPasswordEmail(user.Email, user.ResetPasswordToken); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.MakeError("Failed to send reset password email"))
	}

	return c.SendString("Sent Forgot Password Email")
}
