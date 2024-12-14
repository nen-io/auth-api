package controllers

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

func RefreshToken(c fiber.Ctx) error {

	refreshTojenRequest := new(RefreshTokenRequest)
	err := c.Bind().JSON(&refreshTojenRequest)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	slog.Info("RefreshToken", "refreshToken", refreshTojenRequest.RefreshToken)

	return c.SendString("Refresh Token\n")
}
