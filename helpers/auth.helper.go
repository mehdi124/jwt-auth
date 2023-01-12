package helpers

import (
	"jwt-auth/initializers"
	"jwt-auth/models"
	"jwt-auth/utils/email"
	"jwt-auth/utils/password"
	"gorm.io/gorm"
	"jwt-auth/utils/redis"
	"github.com/thanhpk/randstr"
	"jwt-auth/utils/encode"
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
	user.Email = html.EscapeString(strings.TrimSpace(u.Email))


	err = DB.Select("Email","Password").Create(&user).Error
	if err != nil {
		return "",err
	}

	config, _ := initializers.LoadConfig(".")

	// Generate Verification Code
	code := randstr.String(20)

	verification_code := encode.Encode(code)
	redis.StoreVerificationCode(user.ID,verification_code)

	// ? Send Email
	emailData := utils.EmailData{
		URL:       config.ClientOrigin + "/verifyemail/" + code,
		FirstName: user.Email,
		Subject:   "Your account verification code",
	}

	email.SendEmail(&user, &emailData)

	message := "We sent an email with a verification code to " + user.Email

	return message,nil
}

func Verify(DB *gorm.DB,user *models.User,payload *models.VerifyInput) (string,error) {
	redis.CheckVerificationCode(payload.Email,payload.Code)

	user.EmailVerifiedAt = time.Now().Unix()
	DB.Save(&user)

	Token,err := token.GenerateToken(u.ID)
	if err != nil {
		return "",err
	}

	return Token,nil
}

func Login(DB *gorm.DB,payload *models.LoginInput)(string,error){

	var user models.User
	err := DB.Where("email = ?", payload.Email).Where("email_verified_at NOT NULL").First(&user).Error
	if err != nil{
		return "",err.Error()
	}

	if err := password.VerifyPassword(user.Password,payload.Password); err != nil {
		return "",err
	}

	config, _ := initializers.LoadConfig(".")

	Token, err := utils.GenerateToken(config.TokenExpiresIn, user.ID, config.TokenSecret)
	if err != nil{
		return "",err
	}

	return Token,nil



}


