package controllers

import (
	"log/slog"
	"os"

	"example.com/src/models"
	"example.com/src/services"
	"github.com/gofiber/fiber/v3"
)

func OAuthCallback(c fiber.Ctx) error {
	// This function will handle the OAuth callback logic
	// For now, we will just return a success message

	provider := c.Params("provider")
	if provider == "google" {
		claims, err := services.HandleGoogleCallback(c)

		if err != nil {
			return err
		}

		newUser := models.User{
			Email:     claims.Email,
			UserName:  claims.Email,
			Verified:  true,
			FirstName: "oauth",
			LastName:  "oauth",
			Password:  "oauth",
		}

		exists, err := newUser.EmailExists(newUser.Email)
		if err != nil {
			slog.Error("Error checking if user exists", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.MakeError("Something went wrong"))
		}
		if !exists {
			err := newUser.Create()
			if err != nil {
				slog.Error("Error creating new user", "error", err)
				return c.Status(fiber.StatusInternalServerError).JSON(models.MakeError("Something went wrong"))
			}
		}

		return c.Redirect().To(os.Getenv("UI_DOMAIN") + "/home")
	}

	resp := map[string]any{
		"success": true,
		"message": "OAuth callback successful",
	}

	slog.Info("OAuth Callback", "resp",
		c.Queries())

	return c.JSON(resp)
}
