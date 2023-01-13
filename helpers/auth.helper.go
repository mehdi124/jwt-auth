package helpers

import (
	"html"
	"time"
	"strings"
	"jwt-auth/models"
	"jwt-auth/utils/email"
	"jwt-auth/utils/password"
	"github.com/jinzhu/gorm"
	"jwt-auth/utils/redis"
	"math/rand"
	"jwt-auth/utils/token"
)

func Register(DB *gorm.DB,payload *models.RegisterInput) (string,error) {

	hashedPassword, err := password.HashPassword(payload.Password)
	if err != nil {
		return "",err
	}

	user := models.User{}
	user.Password = string(hashedPassword)
	//remove spaces in email
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))


	err = DB.Select("Email","Password").Create(&user).Error
	if err != nil {
		return "",err
	}


	// Generate Verification Code
	code := string(rand.Intn(999999 - 100000) + 100000)
	redis.StoreVerificationCode(user.ID.String(),code)

	// ? Send Email
	emailData := email.EmailData{
		Code:      code,
		Email: user.Email,
		Subject:   "verification code",
	}

	email.SendEmail(&user, &emailData)

	message := "We sent an email with a verification code to " + user.Email

	return message,nil
}

func Verify(DB *gorm.DB,user *models.User,code string) (string,error) {
	redis.CheckVerificationCode(user.ID.String(),code)

	user.EmailVerifiedAt = time.Now()
	DB.Save(&user)

	Token,err := token.GenerateToken(user.ID.String())
	if err != nil {
		return "",err
	}

	return Token,nil
}

func Login(DB *gorm.DB,payload *models.LoginInput)(string,error){

	var user models.User
	err := DB.Where("email = ?", payload.Email).Where("email_verified_at NOT NULL").First(&user).Error
	if err != nil{
		return "",err
	}

	if err := password.VerifyPassword(user.Password,payload.Password); err != nil {
		return "",err
	}

	Token, err := token.GenerateToken(user.ID.String())
	if err != nil{
		return "",err
	}

	return Token,nil

}


