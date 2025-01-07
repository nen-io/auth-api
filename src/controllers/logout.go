package controllers

import (
	"github.com/gofiber/fiber/v3"
)

func Logout(c fiber.Ctx) error {

	c.ClearCookie("accessToken")
	c.ClearCookie("refreshToken")

	resp := map[string]any{
		"success": true,
		"message": "Logged out",
	}

	return c.JSON(resp)

}
