package main

import (
	"farmsville/backend/database"
	"farmsville/backend/routes"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	database.Connect()
	fmt.Println("Database connected successfully!")
	r := gin.Default()

	routes.SetupStaticRoutes(r)
	routes.SetupAPIRoutes(r)

	// Start the server
	r.Run(":3000")
}
