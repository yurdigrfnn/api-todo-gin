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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()

	// cors middleware
	r.Use(CORSMiddleware())

	// USER ROUTER
	r.POST("api/register", usercontroller.Signup)
	r.POST("api/signin", usercontroller.Signin)
	r.GET("api/validate", middleware.RequireAuth, usercontroller.Validate)
	r.GET("api/logout", usercontroller.Logout)

	//TODO ROUTER
	r.GET("api/todo", middleware.RequireAuth, todocontroller.Index)
	r.POST("api/todo", middleware.RequireAuth, todocontroller.Create)
	r.PUT("api/todo/:id", middleware.RequireAuth, todocontroller.Update)
	r.DELETE("api/todo/:id", middleware.RequireAuth, todocontroller.Destroy)

	r.Run()
}
