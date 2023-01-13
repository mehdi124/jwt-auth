package routes

import (
	"github.com/gin-gonic/gin"
	"jwt-auth/controllers"
	"jwt-auth/middlewares"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {

	router := rg.Group("users")
	router.GET("/me", middlewares.JwtAuthMiddleware(), uc.userController.GetMe)
}
