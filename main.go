package main

import (
	"fmt"
	"log/slog"
	"os"

	"example.com/routes"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		slog.Error("Failed to load .env file")
	}

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))

	app := fiber.New()

	routes.AddRoutes(app)

	app.Listen(port)

	fmt.Println("hello world")
}
