package models

import (
	"time"

	"github.com/alkuinvito/malakh-api/initializers"
)

type Room struct {
	ID        uint
	Name      string
	Picture   string
	Users     []*User `gorm:"many2many:user_rooms"`
	Messages  []Message
	IsPrivate bool `gorm:"not null"`
	CreatedAt time.Time
}

type RoomMaster struct {
	Users []*User `json:"users"`
}

type RoomId struct {
	ID uint `uri:"id" binding:"required"`
}

func (room *Room) PullMessages() error {
	if err := initializers.DB.Table("messages").Where("room_id = ?", room.ID).Find(&(room.Messages)).Error; err != nil {
		return err
	}
	return nil
}

func (room *Room) GetRoomPrivillege(userId string) bool {
	var rows int64
	initializers.DB.Table("user_rooms").Where("room_id = ? AND user_id = ?", room.ID, userId).Count(&rows)

	return rows != 1
}

func InsertRoom(room *Room) (uint, error) {

	result := initializers.DB.Create(&room)
	if result.Error != nil {
		return 0, result.Error
	}

	return room.ID, nil
}

func InsertMessage(message *Message) error {
	err := initializers.DB.Create(&message).Error

	if err != nil {
		return err
	}

	return nil
}

func RoomList(userId string) []Room {
	var roomList []Room

	err := initializers.DB.Table("user_rooms").Where("user_rooms.user_id = ?", userId).Joins("join rooms on rooms.id = user_rooms.room_id").Where("is_private = ?", false).Scan(&roomList).Error
	if err != nil {
		return []Room{}
	}

	return roomList
}
