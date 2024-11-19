package controllers

import (
	"encoding/json"
	"strings"

	"example.com/src/models"
	"example.com/src/services"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func RequestVerify(c fiber.Ctx) error {

	// TODO: use Fiber Bind() to bind the request body to a struct
	body := c.Request().Body()
	// Set a new verification code for the user
	d := json.NewDecoder(strings.NewReader(string(body)))
	d.DisallowUnknownFields()

	request := models.RequestEmail{}
	if err := d.Decode(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.MakeError("Invalid request body"))
	}

	user := models.User{
		Email: request.Email,
	}
	if err := user.GetField("verified"); err != nil {
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
	services.SendVerficationEmail(user.Email, user.VerificationToken)

	return c.SendString("Request Verify")
}
