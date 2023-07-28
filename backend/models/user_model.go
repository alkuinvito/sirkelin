package models

import (
	"time"
)

type User struct {
	ID        string
	Username  string    `json:",omitempty"`
	Fullname  string    `json:",omitempty"`
	Picture   string    `json:",omitempty"`
	Email     string    `gorm:"uniqueIndex;not null" json:",omitempty"`
	Rooms     []*Room   `gorm:"many2many:user_rooms" json:",omitempty"`
	CreatedAt time.Time `json:",omitempty"`
}

type GetByIDParams struct {
	ID string `uri:"id" binding:"required"`
}

type UpdateProfileSchema struct {
	Username string `json:"username" binding:"required,alphanum,lowercase,min=4,max=16"`
	Fullname string `json:"fullname" binding:"required,alphaunicode,min=4,max=32"`
	Picture  string `json:"picture" binding:"required,url"`
}
