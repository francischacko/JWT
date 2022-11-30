package main

import (
	"gin/controller"
	"gin/initializers"
	"gin/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.Loadvariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	r.POST("/signup", controller.Signup)

	r.POST("/login", controller.Login)
	r.GET("/validate", middleware.RequireAuth, controller.Validate)
	r.Run()
}
