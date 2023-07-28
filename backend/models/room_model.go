package models

import "time"

type Room struct {
	ID        string
	Name      string
	Picture   string
	Users     []*User `gorm:"many2many:user_rooms"`
	Messages  []Message
	IsPrivate bool `gorm:"not null"`
	CreatedAt time.Time
}

type CreateRoomParams struct {
	Name      string  `json:"name"`
	Picture   string  `json:"picture"`
	Users     []*User `json:"users"`
	IsPrivate bool    `json:"is_private"`
}

type RoomIDParams struct {
	RoomID string `uri:"id" binding:"required,uuid"`
}

type RoomList struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Picture   string `json:"picture"`
	IsPrivate bool   `json:"is_private"`
}

type UpdateRoomSchema struct {
	Name    string `json:"name" binding:"required,min=3,max=16"`
	Picture string `json:"picture" binding:"url"`
}
