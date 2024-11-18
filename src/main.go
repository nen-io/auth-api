package main

import (
	"fmt"
	"os"

	"example.com/src/routes"
	"example.com/src/services"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {

	// Load envs
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load .env file")
	}

	// initalise services
	services.InitDb()
	services.InitEmailClient()

	app := fiber.New()
	routes.AddRoutes(app)
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	app.Listen(port)

}
