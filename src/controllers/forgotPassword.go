package controllers

import (
	"log/slog"

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

	if err := user.GetField("id"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("Something went wrong"))
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.MakeError("Something went wrong"))
	}

	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("User not found"))
	}

	if err := user.GetField("reset_password_token"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("Something went wrong"))
	}

	newVerificationToken := uuid.New().String()
	if err := user.Update("reset_password_token", newVerificationToken); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.MakeError("Failed to update reset password token"))
	}

	// Send email to user
	if err := services.SendResetPasswordEmail(user.Email, user.ResetPasswordToken, user.ID); err != nil {
		slog.Error("err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.MakeError("Failed to send reset password email"))
	}

	resp := map[string]any{
		"success": true,
		"message": "Sent Forgot Password Email",
	}

	return c.JSON(resp)
}
