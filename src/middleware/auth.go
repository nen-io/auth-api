package middleware

import (
	"fmt"
	"log/slog"
	"os"

	"example.com/src/services"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(c fiber.Ctx) error {
	slog.Info("Auth middleware")

	sess := session.FromContext(c)
	// Check if the provider is google
	provider := sess.Get("provider").(string)
	if provider == "google" {

		tokenInfo, err := services.ValidateGoogleTokens(c)
		if err != nil {
			return err
		}

		slog.Info("Auth middleware", "provider", provider, "tokenInfo", tokenInfo)

		return c.Next()

	}

	accessToken := c.Request().Header.Cookie("accessToken")
	// slog.Info("accessToken", accessToken_c)

	// authHeader := c.Request().Header.Peek("Authorization")
	if accessToken == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Parse the token
	token, err := jwt.Parse(string(accessToken), func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	customClaims := token.Claims.(jwt.MapClaims)

	// Add the user id to the locals for future use
	c.Locals("userId", customClaims["id"])
	c.Locals("email", customClaims["email"])
	c.Locals("username", customClaims["username"])

	return c.Next()
}
