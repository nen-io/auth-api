package services

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func InitSessionStore(app *fiber.App) {
	app.Use(session.New(session.Config{
		// Cookie settings → controls the *session cookie* part:
		CookieHTTPOnly: true,
		CookieSecure:   true,
	}))
}

func GetSession(c fiber.Ctx) (*session.Middleware, error) {
	sess := session.FromContext(c)
	if sess == nil {
		return nil, fiber.ErrUnauthorized
	}
	return sess, nil
}
