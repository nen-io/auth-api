package main

import (
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
	services.DB.AutoMigrate(models.User{})
}
