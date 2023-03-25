package models

import (
	"strings"
	"time"

	"firebase.google.com/go/auth"
	"gorm.io/gorm/clause"

	"github.com/alkuinvito/sirkelin/initializers"
)

type User struct {
	ID        string
	Username  string
	Fullname  string
	Picture   string
	Email     string  `gorm:"uniqueIndex;not null"`
	Rooms     []*Room `gorm:"many2many:user_rooms"`
	CreatedAt time.Time
}

type GetUsersParam struct {
	Fullname string `json:"fullname"`
}

func AuthenticateByIDToken(token *auth.Token) error {
	user := &User{
		ID:       token.Subject,
		Fullname: token.Claims["name"].(string),
		Picture:  token.Claims["picture"].(string),
		Email:    token.Claims["email"].(string),
	}
	return initializers.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&user).Error
}

func GetUsers(fullname string) ([]User, error) {
	var result []User
	err := initializers.DB.Select("id", "fullname", "picture").Where("UPPER(fullname) LIKE ?", "%"+strings.ToUpper(fullname)+"%").Limit(5).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetUserByID(uid string) (*User, error) {
	var result User
	err := initializers.DB.Where("id = ?", uid).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
