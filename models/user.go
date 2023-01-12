package models

import (
	"html"
	"strings"
	"jwt-auth/utils/token"
	"golang.org/x/crypto/bcrypt"
	"errors"
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID  uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Email string `gorm:"size:255;not null;unique" json:"email"`
	Password string `gorm:"size:255;not null;" json:"password"`
	Status bool `gorm:"default:false"`
	EmailVerifiedAt time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt time.Time `gorm:"index;null"`
}

type LoginInput struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterInput struct {
	Email string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type VerifyInput struct {
	Email string `json:"email" binding:"required"`
	Code string `json:"code" binding:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"`
	Photo     string    `json:"photo,omitempty"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


//TODO transfer functions to auth service




func VerifyPassword(password,hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(email , password string) (string,error){

	u := User{}

	err := DB.Model(User{}).Where("email=?",email).Take(&u).Error
	if err != nil{
		return "", err
	}

	err = VerifyPassword(u.Password,password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token,err := token.GenerateToken(u.ID)

	if err != nil {
		return "",err
	}

	return token,nil

}

func GetUserByID(uid uint) (User,error) {

	var u User

	if err := DB.First(&u,uid).Error; err != nil {
		return u,errors.New("User not found!")
	}

	u.PrepareGive()

	return u,nil

}

func (u *User) PrepareGive(){
	u.Password = ""
}


func Drop(){

}