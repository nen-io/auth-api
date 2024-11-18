package routes

import (
	"example.com/src/controllers"
	"github.com/gofiber/fiber/v3"
)

func AddRoutes(app *fiber.App) {
	app.Get("/health", controllers.Health)
	app.Post("/register", controllers.RegisterUser)
	app.Post("/login", controllers.Login)
	app.Get("/verify/:token", controllers.Verify)
	app.Post("/request-verify", controllers.RequestVerify)

}
