package controllers

import (
	"log/slog"

	"example.com/src/models"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func Me(c fiber.Ctx) error {

	sess := session.FromContext(c)

	uIdFromSess := sess.Get("userId")

	if c.Locals("userId") == nil || uIdFromSess == nil || uIdFromSess == "" {
		u := models.User{Email: c.Locals("email").(string)}
		if err := u.GetField("id"); err != nil {
			slog.Error("Failed to get user fields", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.MakeError("Failed to retrieve user information"))
		}
		c.Locals("userId", u.ID)
		sess.Set("userId", u.ID)
	}

	var uId string

	if c.Locals("userId") != nil {
		uId = c.Locals("userId").(string)
	} else {
		uId = uIdFromSess.(string)
	}

	resp := map[string]any{
		"success":  true,
		"message":  "User verified",
		"email":    c.Locals("email"),
		"id":       uId,
		"username": c.Locals("username"),
	}

	return c.JSON(resp)
}
