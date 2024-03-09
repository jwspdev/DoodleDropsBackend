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
	r.Router.GET("/api/validate", middleware.RequireAuth, controllers.Validate)
	r.Router.GET("/api/user/get/current", middleware.RequireAuth, controllers.GetCurrentUser)
	r.Router.POST("/api/user/signup", controllers.Signup)
	r.Router.POST("/api/user/login", controllers.Login)

	//user profile
	r.Router.POST("/api/user/profile/update", middleware.RequireAuth, controllers.UpdateUserProfile)
	//posts
	r.Router.POST("/api/post/create", middleware.RequireAuth, controllers.CreatePost)
	r.Router.POST("/api/post/update", middleware.RequireAuth, controllers.UpdatePost)
	r.Router.POST("/api/post/delete", middleware.RequireAuth, controllers.DeletePost)
	r.Router.GET("/api/post/current", middleware.RequireAuth, controllers.GetCurrentPost)
	r.Router.GET("/api/post/list", middleware.RequireAuth, controllers.ListPosts)
	//comments
	r.Router.POST("/api/comment/create", middleware.RequireAuth, controllers.CreateComment)
	r.Router.POST("/api/comment/delete", middleware.RequireAuth, controllers.DeleteComment)
	//tags
	r.Router.GET("/api/tag/list", middleware.RequireAuth, controllers.ListTags)
	r.Router.POST("/api/tag/user/like", middleware.RequireAuth, controllers.LikeTag)
	r.Router.POST("/api/tag/post/add", middleware.RequireAuth, controllers.AddTagToPost)
}

func (r *Routes) Run() {
	r.Router.Run()
}
