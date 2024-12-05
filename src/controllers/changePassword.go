package controllers

import (
	"example.com/src/models"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

func ChangePassword(c fiber.Ctx) error {
	// Create a struct to hold the request body
	request := new(models.ChangePassword)
	if err := c.Bind().JSON(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("Invalid request body"))
	}

	// Check if the old password is correct
	u := models.User{
		ID: request.Id,
	}

	if err := u.GetField("password"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("Invalid request"))
	}

	if err := u.ComparePassword(request.OldPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("Not allowed to change password"))
	}

	//TODO: Check password strength
	newPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), 14)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.MakeError("Invalid new Password"))
	}

	if err := u.Update("password", string(newPassword)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.MakeError("Failed to update password"))
	}

	return c.SendString("Change Password\n")
}
