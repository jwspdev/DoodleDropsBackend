package routes

import (
	"DoodleDropsBackend/controllers"
	"DoodleDropsBackend/middleware"

	"github.com/gin-gonic/gin"
)

type Routes struct {
	Router *gin.Engine
}

func NewRoutes() *Routes {
	return &Routes{
		Router: gin.Default(),
	}
}

func (r *Routes) SetupRoutes() {
	r.Router.GET("/api/validate", middleware.RequireAuth, controllers.Validate)
	r.Router.GET("/api/user/get", middleware.RequireAuth, controllers.GetCurrentUser)
	r.Router.POST("/api/user/signup", controllers.Signup)
	r.Router.POST("/api/user/login", controllers.Login)

}

func (r *Routes) Run() {
	r.Router.Run()
}
