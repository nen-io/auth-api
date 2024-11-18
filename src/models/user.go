package models

import (
	"errors"
	"log/slog"
	"net/mail"
	"time"

	"example.com/src/database"
	"example.com/src/services"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID                string    `gorm:"type:uuid;primary_key;" json:"id"`
	Email             string    `gorm:"uniqueIndex;not null" json:"email"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Password          string    `json:"password"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	Verified          bool      `json:"verified"`
	VerificationToken string    `json:"verificationToken"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New().String()
	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		slog.Error("Failed to hash password")
		return err
	}
	u.Password = string(pass)
	u.VerificationToken = uuid.New().String()
	return nil
}

func (u *User) ComparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return errors.New("Invalid password")
	}

	return nil
}

func (u *User) GetField(field string) error {

	if u.ID == "" && u.Email == "" {
		return errors.New("User has no ID or Email For GetField")
	}

	if u.ID != "" {
		result := database.DB_Connection.Model(u).Select(field).Where("id = ?", u.ID).First(&u)
		if result.Error != nil {
			return result.Error
		}
		return nil

	} else {
		result := database.DB_Connection.Model(u).Select(field).Where("email = ?", u.Email).First(&u)
		if result.Error != nil {
			return result.Error
		}
		return nil
	}

}

func (u *User) GetFields(fields []string) error {

	if u.ID == "" && u.Email == "" {
		return errors.New("User has no ID or Email For GetField")
	}

	if u.ID != "" {
		result := database.DB_Connection.Model(u).Select(fields).Where("id = ?", u.ID).First(&u)
		if result.Error != nil {
			return result.Error
		}
		return nil
	} else {
		result := database.DB_Connection.Model(u).Select(fields).Where("email = ?", u.Email).First(&u)
		if result.Error != nil {
			return result.Error
		}
		return nil
	}
}

func (u *User) Update(field string, value any) error {

	if u.ID == "" && u.Email == "" {
		return errors.New("User has no ID or Email For Update")
	}

	if u.ID != "" {
		result := database.DB_Connection.Model(u).Update(field, value).Where("id = ?", u.ID)
		if result.Error != nil {
			return result.Error
		}
		return nil
	} else {
		result := database.DB_Connection.Model(u).Update(field, value).Where("email = ?", u.Email)
		if result.Error != nil {
			return result.Error
		}
		return nil
	}
}

func (u *User) UpdateFields(fields map[string]any) error {

	if u.ID == "" && u.Email == "" {
		return errors.New("User has no ID or Email For GetField")
	}

	if u.ID != "" {
		result := database.DB_Connection.Model(u).Updates(fields).Where("id", u.ID)
		if result.Error != nil {
			return result.Error
		}
		return nil
	} else {
		result := database.DB_Connection.Model(u).Updates(fields).Where("email", u.Email)
		if result.Error != nil {
			return result.Error
		}
		return nil

	}

}

func (u *User) SendVerificationEmail() error {
	// send email
	return nil
}

func (u *User) ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) EmailExists(email string) bool {
	if err := database.DB_Connection.First(&User{Email: email}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}

	return true
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
	if u.EmailExists(u.Email) {
		return errors.New("Email already in use or something went wrong checking")
	}

	if len(u.Password) < 6 { // validate password length
		return errors.New("Password too short")
	}

	return nil

}

func (u *User) Create() error {
	result := database.DB_Connection.Create(u)
	if result.Error != nil {
		return result.Error
	}
	if err := services.SendVerficationEmail(u.Email, u.VerificationToken); err != nil {
		return err
	}
	return nil
}
