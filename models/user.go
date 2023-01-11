package models

import (
	"html"
	"strings"
	"jwt-auth/utils/token"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"errors"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null;" json:"password"`
}

func (u *User) SaveUser() (*User, error) {

	var err error

	err = u.BeforeSave()
	if err != nil{
		return &User{},err
	}

	ConnectDatabase()

	err = DB.Select("Username","Password").Create(&u).Error
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
	//remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	return nil

}

func VerifyPassword(password,hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username , password string) (string,error){

	u := User{}

	err := DB.Model(User{}).Where("username=?",username).Take(&u).Error
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
