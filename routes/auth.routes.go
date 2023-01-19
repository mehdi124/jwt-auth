package routes

import (
	"github.com/gin-gonic/gin"
	"jwt-auth/controllers"
	"jwt-auth/middlewares"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (rc *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("/auth")

	router.GET("/test", rc.authController.Test)
	router.POST("/register", rc.authController.Register)
	router.POST("/login", rc.authController.Login)
	router.DELETE("/logout", middlewares.JwtAuthMiddleware(), rc.authController.Logout)
	router.POST("/verify", rc.authController.Verify)
}
