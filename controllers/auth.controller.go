package controllers

import (
	"jwt-auth/helpers"
	"jwt-auth/utils/token"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"jwt-auth/initializers"
	"jwt-auth/models"
	"jwt-auth/utils/password"
	"jwt-auth/utils/redis"
	"jwt-auth/utils/email"
	"jwt-auth/utils/encode"
	"gorm.io/gorm"
)

type AuthController struct{
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{ DB:DB }
}

func (ac *AuthController) Register(ctx *gin.Context){

	var payload models.RegisterInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if payload.Password != payload.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	message,err := helpers.Register(ac.DB,&payload)
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

	var payload models.VerifyInput
	var user models.User

	result := ac.DB.Where("email = ?", payload.Email).Where("email_verified_at IS NULL").First(&user)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid verification code or user doesn't exists"})
		return
	}

	redis.CheckVerificationCode(payload.Email,payload.Code)

	user.EmailVerifiedAt = time.Now().Unix()
	ac.DB.Save(&user)

	t,err := token.GenerateToken(u.ID)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "token": t})
}

