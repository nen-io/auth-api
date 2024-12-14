package middleware

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(c fiber.Ctx) error {
	slog.Info("Auth middleware")

	// headers := c.GetReqHeaders()

	authHeader := c.Request().Header.Peek("Authorization")
	if authHeader == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	accessToken := strings.Split(string(authHeader), " ")[1]
	slog.Info(accessToken)
	if accessToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Parse the token
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
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

	fmt.Printf("%v", customClaims)

	// Add the user id to the locals for future use
	c.Locals("userId", customClaims["id"])
	c.Locals("email", customClaims["email"])

	return c.Next()
}
