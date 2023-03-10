package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yurdigrfnn/api-todo-auth/controllers/todocontroller"
	"github.com/yurdigrfnn/api-todo-auth/controllers/usercontroller"
	"github.com/yurdigrfnn/api-todo-auth/initializers"
	"github.com/yurdigrfnn/api-todo-auth/middleware"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDatabase()
}

func main() {
	r := gin.Default()
	// USER ROUTER
	r.POST("api/register", usercontroller.Signup)
	r.POST("api/signin", usercontroller.Signin)
	r.GET("api/validate", middleware.RequireAuth, usercontroller.Validate)

	//TODO ROUTER
	r.GET("api/todo", middleware.RequireAuth, todocontroller.Index)
	r.POST("api/todo", middleware.RequireAuth, todocontroller.Create)
	r.PUT("api/todo/:id", middleware.RequireAuth, todocontroller.Update)
	r.DELETE("api/todo/:id", middleware.RequireAuth, todocontroller.Destroy)

	r.Run()
}
