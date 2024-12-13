package main

import (
	"log/slog"

	"example.com/src/models"
	"example.com/src/services"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load .env file")
	}
	services.InitDb()
	err = services.DB.AutoMigrate(models.User{})
	if err != nil {
		panic("Failed to migrate database")
	}
	slog.Info("Database migration completed...")
}
