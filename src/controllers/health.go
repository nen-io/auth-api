package controllers

import "github.com/gofiber/fiber/v3"

func Health(c fiber.Ctx) error {

	resp := map[string]string{
		"status": "ok",
	}

	return c.JSON(&resp)
}
