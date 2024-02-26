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
	//user
	r.Router.GET("/api/user/validate", middleware.RequireAuth, controllers.Validate)
	r.Router.GET("/api/user/get", middleware.RequireAuth, controllers.GetCurrentUser)
	r.Router.POST("/api/user/signup", controllers.Signup)
	r.Router.POST("/api/user/login", controllers.Login)

	//user profile
	r.Router.POST("/api/user/profile/update", middleware.RequireAuth, controllers.UpdateUserProfile)
}

func (r *Routes) Run(address string) {
	r.Router.Run(address)
}
