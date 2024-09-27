package models

import (
	"github.com/sid/FileSystem/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Name            string `json:"name"`
	Age             string `json:"age"`
	Email           string `gorm:"primaryKey" json:"username"`
	EncryptionToken string `json:"encryption_token`
	Password        string `json:"password"`
}

var db gorm.DB

func init() {
	db = *config.Connect()
	db.AutoMigrate(&User{})
}
func (u *User) CreateUser() *User {
	hashPass, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hashPass)
	db.Create(&u)
	return u
}

func GetByEmail(email string) (*User, *gorm.DB) {
	var user User
	db := db.Where("Email=?", email).Find(&user)
	if db != nil {
		return &user, db
	}
	return nil, nil
}
