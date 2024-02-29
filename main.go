package main

import (
	"DoodleDropsBackend/initializers"
	"DoodleDropsBackend/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	// Create an instance of Routes
	routes := routes.NewRoutes()

	// Setup routes
	routes.SetupRoutes()

	// Run the Gin server
	routes.Run()
	r.Run()
}
