package routes

import (
	"example.com/controllers"
	"github.com/gofiber/fiber/v3"
)

func AddRoutes(app *fiber.App) {
	app.Get("/health", controllers.Health)
	app.Post("/register", controllers.RegisterUser)

}
