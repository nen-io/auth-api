package main

import (
	"fmt"
	"os"

	"example.com/src/routes"
	"example.com/src/services"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

type structValidator struct {
	validate *validator.Validate
}

// Validate method implementation
func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func main() {

	// Load envs
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load .env file")
	}

	// initalise services
	services.InitDb()
	services.InitEmailClient()
	services.InitSession()

	// Add validator and fast json encoder/decoder
	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{validate: validator.New(validator.WithRequiredStructEnabled())},
		JSONEncoder:     json.Marshal,
		JSONDecoder:     json.Unmarshal,
	})
	routes.AddRoutes(app)
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	app.Listen(port)

}
