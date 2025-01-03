package models

import (
	"errors"
	"log/slog"
	"net/mail"
	"time"

	"example.com/src/services"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID                 string    `gorm:"type:uuid;primary_key;" json:"id"`
	Email              string    `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	UserName           string    `gorm:"uniqueIndex;not null" json:"userName" validate:"required"`
	FirstName          string    `json:"firstName" validate:"required"`
	LastName           string    `json:"lastName" validate:"required"`
	Password           string    `json:"password" validate:"required"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
	Verified           bool      `json:"verified"`
	VerificationToken  string    `json:"verificationToken"`
	ResetPasswordToken string    `json:"resetPasswordToken"`
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
		result := services.DB.Model(u).Select(field).Where("id = ?", u.ID).First(&u)
		if result.Error != nil {
			return result.Error
		}
		return nil

	} else if u.Email != "" {
		result := services.DB.Model(u).Select(field).Where("email = ?", u.Email).First(&u)
		if result.Error != nil {
			return result.Error
		}
		return nil
	} else {

		result := services.DB.Model(u).Select(field).Where("user_name = ?", u.Email).First(&u)
		if result.Error != nil {
			return result.Error
		}
	}

	return errors.New("Failed to find user")

}

func (u *User) GetFields(fields []string) error {

	if u.ID == "" && u.Email == "" && u.UserName == "" {
		return errors.New("User has no ID or Email For GetField")
	}

	if u.ID != "" {
		result := services.DB.Model(u).Select(fields).Where("id = ?", u.ID).First(&u)
		if result.Error != nil {
			return result.Error
		}
		return nil
	} else if u.Email != "" {
		result := services.DB.Model(u).Select(fields).Where("email = ?", u.Email).First(&u)
		if result.Error != nil {
			return result.Error
		}
		return nil
	} else {
		result := services.DB.Model(u).Select(fields).Where("user_name = ?", u.UserName).First(&u)
		if result.Error != nil {
			return result.Error
		}
	}

	return errors.New("Failed to find user")
}

func (u *User) Update(field string, value any) error {

	if u.ID == "" && u.Email == "" {
		return errors.New("User has no ID or Email For Update")
	}

	if u.ID != "" {
		result := services.DB.Model(u).Update(field, value).Where("id = ?", u.ID)
		if result.Error != nil {
			return result.Error
		}
		return nil
	} else {
		result := services.DB.Model(u).Update(field, value).Where("email = ?", u.Email)
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
		result := services.DB.Model(u).Updates(fields).Where("id", u.ID)
		if result.Error != nil {
			return result.Error
		}
		return nil
	} else {
		result := services.DB.Model(u).Updates(fields).Where("email", u.Email)
		if result.Error != nil {
			return result.Error
		}
		return nil

	}

}

func (u *User) ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) EmailExists(email string) (bool, error) {
	if err := services.DB.First(&User{Email: email}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return true, err

	}

	return true, nil
}

func (u *User) Validate() error {
	// check if user already exists
	exists, err := u.EmailExists(u.Email)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("Email already in use or something went wrong checking")
	}

	//TODO: Check password strength
	if len(u.Password) < 6 { // validate password length
		return errors.New("Password too short")
	}

	return nil

}

func (u *User) Create() error {
	result := services.DB.Create(u)
	if result.Error != nil {
		return result.Error
	}
	if err := services.SendVerficationEmail(u.Email, u.VerificationToken, u.ID); err != nil {
		return err
	}
	return nil
}
