package main

import (
	"fmt"
	"os"

	"example.com/src/routes"
	"example.com/src/services"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"

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
	services.InitGoogleOAuth()

	uiDomain := os.Getenv("UI_DOMAIN")

	// Add validator and fast json encoder/decoder
	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{validate: validator.New(validator.WithRequiredStructEnabled())},
		JSONEncoder:     json.Marshal,
		JSONDecoder:     json.Unmarshal,
	})

	services.InitSessionStore(app)

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     []string{uiDomain, "http://192.168.0.226:3001"},
	}))
	routes.AddRoutes(app)
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	app.Listen(port)

}
