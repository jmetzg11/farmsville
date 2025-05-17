package main

import (
	"farmsville/backend/database"
	"farmsville/backend/routes"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found. Using environment variables.")
	}

	database.Connect()
	fmt.Println("Database connected successfully!")
	r := gin.Default()

	r.Static("/photos", "./data/photos")
	routes.SetupStaticRoutes(r)
	routes.SetupAPIRoutes(r)

	// Start the server
	r.Run(":3000")
}
