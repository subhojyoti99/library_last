// model/user_model.go

package model

import (
	"Reference/database"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID            uint   `gorm:"primary_key;auto_increment"`
	Name          string `json:"name"`
	Email         string `gorm:"required" json:"email"`
	ContactNumber string `json:"contact_number"`
	Password      string `json:"password"`
	Role          string `json:"role"`
	RoleID        uint   `gorm:"not null;DEFAULT:3" json:"role_id"`
	LibID         uint   `json:"lib_id"`
	ValidKey      string `json:"valid_key"`
}

func (u *User) Save() (*User, error) {
	err := database.DB.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, err
}

// Generate encrypted password
func (user *User) BeforeSave(*gorm.DB) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	return nil
}

// Get all users
func GetUsers(User *[]User) (err error) {
	err = database.DB.Find(User).Error
	if err != nil {
		return err
	}
	return nil
}

// Get user by email
func GetUserByEmail(email string) (User, error) {
	var user User
	err := database.DB.Where("email=?", email).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// Validate user password
func (user *User) ValidateUserPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

// Get user by id
func GetUserById(id uint) (User, error) {
	var user User
	err := database.DB.Where("id=?", id).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// Get user by id
func GetUser(User *User, id int) (err error) {
	err = database.DB.Where("id = ?", id).First(User).Error
	if err != nil {
		return err
	}
	return nil
}

// Update user
func UpdateUser(User *User) (err error) {
	err = database.DB.Omit("password").Updates(User).Error
	if err != nil {
		return err
	}
	return nil
}
