package controllers

import (
	"log/slog"

	"example.com/src/services"
	"github.com/gofiber/fiber/v3"
)

func OAuthCallback(c fiber.Ctx) error {
	// This function will handle the OAuth callback logic
	// For now, we will just return a success message

	provider := c.Params("provider")
	if provider == "google" {
		return services.HandleGoogleCallback(c)
	}

	resp := map[string]any{
		"success": true,
		"message": "OAuth callback successful",
	}

	slog.Info("OAuth Callback", "resp",
		c.Queries())

	return c.JSON(resp)
}
