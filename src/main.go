package main

import (
	"fmt"
	"os"

	"example.com/src/database"
	"example.com/src/routes"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic("Failed to load .env file")
	}

	database.Connect()

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))

	app := fiber.New()

	routes.AddRoutes(app)

	app.Listen(port)

}
