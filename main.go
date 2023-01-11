package main

import (
	"jwt-auth/controllers"
	"jwt-auth/models"
	"jwt-auth/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {

	models.ConnectDatabase()

	r := gin.Default()

	public := r.Group("/api")

	public.POST( "/register",controllers.Register )
	public.POST( "/login",controllers.Login )

	protected := r.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user",controllers.CurrentUser)

	r.Run(":8000")
}

