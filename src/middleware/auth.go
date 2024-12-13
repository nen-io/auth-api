package middleware

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

func Auth(c fiber.Ctx) error {
	slog.Info("Auth middleware")

	// headers := c.GetReqHeaders()

	key := c.Request().Header.Peek("Authorization")

	fmt.Printf("%s \n\n", key)
	return c.Next()
}
