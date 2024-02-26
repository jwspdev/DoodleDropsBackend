package main

import (
	"DoodleDropsBackend/controllers"
	"DoodleDropsBackend/initializers"
	"DoodleDropsBackend/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.GET("/api/validate", middleware.RequireAuth, controllers.Validate)
	r.POST("/api/signup", controllers.Signup)
	r.POST("/api/login", controllers.Login)
	r.Run()
}
