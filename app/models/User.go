package models

import (
	"errors"
	"fmt"
	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

type User struct {
	gorm.Model
	Name string `gorm:"size:100;not null"`
	Email string `gorm:"size:100;not null"`
	Password string `gorm:"size:255;not null"`
}

// HashPassword
func HashPassword(password string) (string) {
	hash, err :=  bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

// CompareHashPassword
func CheckHashPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("Invalid Credentials")
	}
	return nil
}

// BeforeSave hash the user password
func (user *User) BeforeSave() {
	password := strings.TrimSpace(user.Password)
	user.Password = HashPassword(password)
}

// Prepare strips user input of any whitespaces
func (user *User) Prepare() {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)
}

// Validate the user input
func (user *User) Validate(action string) error {
	switch(strings.ToUpper(action)) {
	case "REGISTER":
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			errors.New("Email is required and must be valid")
		}
		if user.Password == "" {
			errors.New("Password is required")
		}
		return nil
	case "LOGIN":
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Email is required and must be valid")
		}
		if user.Password == "" {
			return errors.New("Password is required")
		}
		return nil
	default:
		return nil
	}
}

// SaveUser to database
func (user *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error

	err = db.Debug().Create(&user).Error
	if err != nil {
		fmt.Println("Something wrong: ", err)
		return &User{}, err
	}
	return user, err
}

// GetUser get one user
func (user *User) GetUser(db *gorm.DB) (*User, error) {
	account := &User{}
	if err := db.Debug().Table("users").Where("email = ?", user.Email).First(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}