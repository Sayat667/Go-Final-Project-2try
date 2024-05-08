package main

import (
	"log"
	"gofp/controllers"
	"gofp/models"
	"gofp/routes"
	"gofp/utils"
	// "gofp/controllers"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

func main() {
	utils.Setenv()
	// Create a new gin instance
	r := gin.Default()

	// Load .env file and Create a new connection to the database
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	config := models.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	// Initialize DB
	db := models.InitDB(config)
	SetDB(db)
	// Load the routes
	routes.AuthRoutes(r)

	// Run the server
	r.Run(":8080")
}

func SetDB(db *gorm.DB) {
	controllers.DB = db
	routes.DB = db
}
