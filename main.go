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

	routes := routes.NewRoutes()

	// Setup routes
	routes.SetupRoutes()

	//Run routes
	routes.Run()
	// Run the Gin server
	r.Run()
}
