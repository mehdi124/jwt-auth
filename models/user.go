package models

import (
	"time"
)

type User struct {
	ID        uint    `json:"id" gorm:"autoIncrement; primaryKey"`
	Email   string    `gorm:"unique;index;not null"`
	Password  string    `gorm:"not null"`
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
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type VerifyInput struct {
	Email string `json:"email" binding:"required,email"`
	Code string `json:"code" binding:"required"`
}

type UserResponse struct {
	ID        uint `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"`
	Photo     string    `json:"photo,omitempty"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}



//func LoginCheck(email , password string) (string,error){
//
//	u := User{}
//
//	err := DB.Model(User{}).Where("email=?",email).Take(&u).Error
//	if err != nil{
//		return "", err
//	}
//
//	err = VerifyPassword(u.Password,password)
//
//	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
//		return "", err
//	}
//
//	token,err := token.GenerateToken(u.ID)
//
//	if err != nil {
//		return "",err
//	}
//
//	return token,nil
//
//}
//
//func GetUserByID(uid uint) (User,error) {
//
//	var u User
//
//	if err := DB.First(&u,uid).Error; err != nil {
//		return u,errors.New("User not found!")
//	}
//
//	u.PrepareGive()
//
//	return u,nil
//
//}
//
//func (u *User) PrepareGive(){
//	u.Password = ""
//}
//
//
//func Drop(){
//
//}