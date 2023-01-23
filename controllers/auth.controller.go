package controllers

import (
	"jwt-auth/helpers"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"jwt-auth/models"
	"github.com/jinzhu/gorm"
)

type AuthController struct{
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{ DB:DB }
}

//func (ac *AuthController) Test(ctx *gin.Context){
//
//	helpers.Test()
//	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "test"})
//
//}

func (ac *AuthController) Register(ctx *gin.Context){

	var payload *models.RegisterInput


	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if payload.Password != payload.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}


	message,err := helpers.Register(ac.DB,payload)
	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": message})

}

func (ac *AuthController) Verify(ctx *gin.Context) {

	var payload *models.VerifyInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user := models.User{}
	err := ac.DB.Where("email = ? AND email_verified_at IS NULL",payload.Email).First(&user).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid verification code or user doesn't exists"})
		return
	}

	//log.Println(user,"ssss")
	token,err := helpers.Verify(ac.DB,&user,payload.Code)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened"})
		return
	}

	//TODO expire time is not dynamic
	ctx.SetCookie("token", token, 60*60, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "token": token})
}


func (ac *AuthController) Login(ctx *gin.Context) {

	var payload *models.LoginInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}


	token,err := helpers.Login(ac.DB,payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err})
		return
	}

	//TODO expire time is not dynamic
	ctx.SetCookie("token", token, 60*60, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "token": token})
}

func (ac *AuthController) Logout(ctx *gin.Context) {

	ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})

}


