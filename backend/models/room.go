package models

import (
	"github.com/google/uuid"
	"time"

	"github.com/alkuinvito/sirkelin/initializers"
)

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

func (room *Room) PullMessages() error {
	return initializers.DB.Table("messages").Where("room_id = ?", room.ID).Find(&(room.Messages)).Error
}

func (room *Room) GetRoomPrivilege(uid string) bool {
	var rows int64
	initializers.DB.Table("user_rooms").Where("room_id = ? AND user_id = ?", room.ID, uid).Count(&rows)

	return rows != 1
}

func InsertRoom(room *CreateRoomParams) (string, error) {
	id := uuid.New().String()
	err := initializers.DB.Create(&Room{
		ID:        id,
		Name:      room.Name,
		Picture:   room.Picture,
		Users:     room.Users,
		Messages:  nil,
		IsPrivate: room.IsPrivate,
		CreatedAt: time.Time{},
	}).Error
	if err != nil {
		return "", err
	}

	return id, nil
}

func RoomList(uid string) []Room {
	var result []Room

	err := initializers.DB.Table("user_rooms").Where("user_rooms.user_id = ?", uid).Joins("join rooms on rooms.id = user_rooms.room_id").Where("is_private = ?", false).Scan(&result).Error
	if err != nil {
		return []Room{}
	}

	return result
}
