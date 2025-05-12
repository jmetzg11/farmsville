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

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Connect()
	fmt.Println("Database connected successfully!")
	r := gin.Default()

	routes.SetupStaticRoutes(r)
	routes.SetupAPIRoutes(r)

	// Start the server
	r.Run(":3000")
}
