package models

import (
	"errors"
	"log/slog"
	"net/mail"

	"example.com/src/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        string `gorm:"type:uuid;primary_key;" json:"id"`
	Email     string `gorm:"uniqueIndex;not null" json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New().String()
	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		slog.Error("Failed to hash password")
		return err
	}
	u.Password = string(pass)
	return nil
}

func (u *User) ComparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return errors.New("Invalid password")
	}

	return nil
}

func (u *User) GetByEmail(email string) error {
	result := database.DB_Connection.Where("email = ?", email).First(&u)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *User) ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Validate() error {
	if u.FirstName == "" || u.LastName == "" || u.Email == "" || u.Password == "" {
		return errors.New("Empty fields")
	}

	// validate email
	if err := u.ValidateEmail(u.Email); err != nil {
		return err
	}

	// check if user already exists
	eu := &User{}
	if err := eu.GetByEmail(u.Email); eu.Email != "" || err == nil {
		return errors.New("Email already in use or something went wrong checking")
	}

	if len(u.Password) < 6 { // validate password length
		return errors.New("Password too short")
	}

	return nil

}

func (u *User) Save() error {
	result := database.DB_Connection.Create(u)
	if result.Error != nil {
		return result.Error
	}
	slog.Info("create", "result", result.Statement)
	return nil
}
