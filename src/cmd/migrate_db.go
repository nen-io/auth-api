package main

import (
	"example.com/src/database"
	"example.com/src/models"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load .env file")
	}
	database.Connect()
	database.DB_Connection.AutoMigrate(models.User{})
}
