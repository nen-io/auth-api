package controllers

import (
	"encoding/json"
	"strings"

	"example.com/src/models"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

func ChangePassword(c fiber.Ctx) error {
	// This function will be used to change a user's password

	body := c.Request().Body()
	// Set a new verification code for the user
	d := json.NewDecoder(strings.NewReader(string(body)))
	d.DisallowUnknownFields()

	// Create a struct to hold the request body
	request := models.ChangePassword{}
	if err := d.Decode(&request); err != nil {
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
