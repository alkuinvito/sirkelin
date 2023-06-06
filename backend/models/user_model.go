package models

import (
	"time"
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

type GetByIDParams struct {
	ID string `uri:"id" binding:"required"`
}
