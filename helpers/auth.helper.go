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
	"strconv"
	"jwt-auth/utils/token"
	"errors"
)

func Test(){
	emailData := email.EmailData{
		Code:      "1234",
		Email: "admin@admin.com",
		Subject:   "verification code",
	}
	email.RunSendEmailJob(emailData,"register")
}

func Register(DB *gorm.DB,payload *models.RegisterInput) (string,error) {

	user := models.User{}
	err := DB.Where("email = ?",payload.Email).First(&user).Error

	if !user.EmailVerifiedAt.IsZero() {
		return "", errors.New("user not exist")
	}


	hashedPassword, Err := password.HashPassword(payload.Password)
	if Err != nil {
		return "",Err
	}

	//record not found
	if err != nil{

		user.Password = string(hashedPassword)
		//remove spaces in email
		user.Email = html.EscapeString(strings.TrimSpace(payload.Email))

		err = DB.Select("Email","Password").Create(&user).Error
		if err != nil {
			return "",err
		}

	}else{

		user.Password = string(hashedPassword)

		err = DB.Select("Password").Save(&user).Error
		if err != nil {
			return "",err
		}

	}


	// Generate Verification Code
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(999999 - 100000) + 100000
	Code := strconv.Itoa(code)
	redis.StoreVerificationCode(user.ID,Code)

	// ? Send Email
	emailData := email.EmailData{
		Code:      Code,
		Email: user.Email,
		Subject:   "verification code",
	}

	email.SendEmail(&user, &emailData)

	message := "We sent an email with a verification code to " + user.Email

	return message,nil
}

func Verify(DB *gorm.DB,user *models.User,code string) (string,error) {
	redis.CheckVerificationCode(user.ID,code)

	tx := DB.Begin()

	user.EmailVerifiedAt = time.Now()
	err := tx.Select("EmailVerifiedAt").Save(&user).Error
	if err != nil {
		tx.Rollback()
		return "",err
	}

	Token,err := token.GenerateToken(user.ID)
	if err != nil {
		tx.Rollback()
		return "",err
	}

	tx.Commit()
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

	Token, err := token.GenerateToken(user.ID)
	if err != nil{
		return "",err
	}

	return Token,nil

}


