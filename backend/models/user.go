package models

import (
	"firebase.google.com/go/auth"
	"time"

	"github.com/alkuinvito/sirkelin/initializers"
	"gorm.io/gorm/clause"
)

type User struct {
	ID        string
	Username  string `gorm:"uniqueIndex;not null"`
	Fullname  string
	Picture   string
	Email     string  `gorm:"uniqueIndex;not null"`
	Rooms     []*Room `gorm:"many2many:user_rooms"`
	CreatedAt time.Time
}

func AuthenticateByIDToken(token *auth.Token) {
	user := &User{
		ID:       token.Subject,
		Fullname: token.Claims["name"].(string),
		Picture:  token.Claims["picture"].(string),
		Email:    token.Claims["email"].(string),
	}
	initializers.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)
}

func GetUserByID(uid string) (*User, error) {
	var result User
	err := initializers.DB.Where("id = ?", uid).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
