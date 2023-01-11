package models

import (
	"html"
	"strings"
	"jwt-auth/utils/token"
	"golang.org/x/crypto/bcrypt"
	"errors"
	"time"
)

type User struct {
	ID uint `gorm:"primaryKey;autoIncrement"`
	Email string `gorm:"size:255;not null;unique" json:"email"`
	Password string `gorm:"size:255;not null;" json:"password"`
	Status bool `gorm:"default:false"`
	EmailVerifiedAt time.Time `gorm:"index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt time.Time `gorm:"index;null"`
}

func (u *User) SaveUser() (*User, error) {

	var err error

	err = u.BeforeSave()
	if err != nil{
		return &User{},err
	}

	db := ConnectDatabase()

	err = db.Select("Email","Password").Create(&u).Error
	if err != nil {
		return &User{}, err
	}

	return u, nil
}

func (u *User) BeforeSave() error {

	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password),bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	//remove spaces in email
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	return nil

}

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