package models

import (
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

func (user *User) UserAuthenticate() {
	initializers.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)
}
