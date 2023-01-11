package cli

import (
	"github.com/gin-gonic/gin"
	"jwt-auth/controllers"
	"jwt-auth/middlewares"
)

func (cli *CLI)Serve(port string){

	r := gin.Default()

	public := r.Group("/api")

	public.POST( "/register",controllers.Register )
	public.POST( "/login",controllers.Login )

	protected := r.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user",controllers.CurrentUser)

	r.Run(":" + port)

}
