package controllers

import (
	"log/slog"

	"example.com/src/models"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

func ForgotPasswordChange(c fiber.Ctx) error {
	request := new(models.ForgotPasswordChange)

	if err := c.Bind().JSON(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("Invalid request body"))
	}

	user := models.User{
		ID: request.ID,
	}

	if err := user.GetField("reset_password_token"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("Something went wrong"))
	}

	if user.ResetPasswordToken != request.Token {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("Invalid reset password token"))
	}

	// create a hash of the new password
	pass, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), 14)
	if err != nil {
		slog.Error("Failed to hash password")
		return err
	}

	// update the user's password in the database with the new hash
	if err := user.Update("password", string(pass)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.MakeError("Failed to update password"))
	}

	resp := map[string]any{
		"success": true,
		"message": "Password changed",
	}

	return c.JSON(resp)
}
