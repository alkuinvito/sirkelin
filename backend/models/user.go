package models

import (
	"fmt"
	"html"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/alkuinvito/malakh-api/initializers"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID uint
	Username     string  `gorm:"uniqueIndex;not null"`
	Fullname     string
	Picture string
	Email    string  `gorm:"uniqueIndex;not null"`
	Password string  `gorm:"not null"`
	Rooms    []*Room `gorm:"many2many:user_rooms"`
	CreatedAt time.Time
}

func (user *User) SaveUser() (*User, error) {
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	is_alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(user.Username)
	if !is_alphanumeric || (len(user.Username) < 4 || len(user.Username) > 12) {
		return &User{}, fmt.Errorf("username must be an alphanumeric")
	}

	if len(user.Password) < 8 || len(user.Password) > 32 {
		return &User{}, fmt.Errorf("password length must be between 8 and 32 characters")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return &User{}, err
	}
	user.Password = string(hashed)

	if _, err = mail.ParseAddress(user.Email); err != nil {
		return &User{}, err
	}

	if err = initializers.DB.Create(&user).Error; err != nil {
		return &User{}, fmt.Errorf("username and/or email already exists")
	}

	return user, nil
}

func ValidatePassword(password, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}

func (user *User) ValidateLogin() (*User, error) {
	var err error
	result := User{}

	if err = initializers.DB.Model(&User{}).Where("username = ?", user.Username).First(&result).Error; err != nil {
		return &User{}, err
	}

	err = ValidatePassword(user.Password, result.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return &User{}, err
	}

	return &result, nil
}
