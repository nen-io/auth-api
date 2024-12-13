package routes

import (
	"example.com/src/controllers"
	"example.com/src/middleware"
	"github.com/gofiber/fiber/v3"
)

func AddRoutes(app *fiber.App) {
	app.Get("/health", controllers.Health)

	app.Post("/login", controllers.Login)

	app.Post("/register", controllers.RegisterUser)
	app.Post("/request-email-verification", controllers.RequestEmailVerificationCode)
	app.Get("/verify/:token/:id", controllers.Verify)

	app.Post("/forgot-password", controllers.ForgotPassword)
	app.Post("/forgot-password-change", controllers.ForgotPasswordChange)

	app.Post("/change-password", controllers.ChangePassword, middleware.Auth)
	app.Get("/auth", controllers.Restricted, middleware.Auth)
}
