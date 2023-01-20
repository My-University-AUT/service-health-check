package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email          string `gorm:"unique" validate:"required"`
	HashedPassword string
	Links          []Link
}

type UserToLogin struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserToRegister struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type UserResponse struct {
	ID    uint
	Email string `json:"email,omitempty"`
	Links []Link
}

func (user *User) CreateUser(db *gorm.DB) error {
	hashedPassword, err := HashPassword(user.HashedPassword)
	if err != nil {
		return err
	}

	user.HashedPassword = hashedPassword
	result := db.Create(user)

	return result.Error
}

func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	var user User
	result := db.Model(&User{}).Where(&User{Email: email}).Preload("Links").First(&user)
	return &user, result.Error
}

func GetUserByEmailPassword(db *gorm.DB, username, password string) (*User, error) {
	user, err := GetUserByEmail(db, username)

	if err != nil {
		return nil, err
	}

	err = VerifyPassword(user.HashedPassword, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("could not hash password %w", err)
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword string, candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
}
