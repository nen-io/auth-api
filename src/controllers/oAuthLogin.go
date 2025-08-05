package controllers

import (
	"example.com/src/services"
	"github.com/gofiber/fiber/v3"
)

func OAuthLogin(c fiber.Ctx) error {

	provider := c.Params("provider") // e.g., "google", "github", etc.
	if provider == "google" {
		return services.HandleGoogleAuth(c)
	}

	return c.Status(fiber.StatusBadRequest).JSON(map[string]string{
		"error": "Unsupported provider",
	})

}
