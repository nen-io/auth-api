package controllers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func Logout(c fiber.Ctx) error {

	c.ClearCookie("accessToken")
	c.ClearCookie("refreshToken")
	sess := session.FromContext(c)
	sess.Destroy()

	resp := map[string]any{
		"success": true,
		"message": "Logged out",
	}

	return c.JSON(resp)

}
