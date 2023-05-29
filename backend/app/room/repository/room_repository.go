package repository

import (
	"sirkelin/backend/models"

	"gorm.io/gorm"
)

type RoomRepository struct {
}

type IRoomRepository interface {
	Count(db *gorm.DB, room *models.Room, uid string) (int, error)
	Create(db *gorm.DB, room *models.Room) error
	GetByID(db *gorm.DB, uid string) ([]models.Room, error)
	Peek(db *gorm.DB, room *models.Room) error
}

func NewRoomRepository() *RoomRepository {
	return &RoomRepository{}
}

func (repo *RoomRepository) Count(db *gorm.DB, room *models.Room, uid string) (int, error) {
	var rows int64
	err := db.Table("user_rooms").Where("room_id = ? AND user_id = ?", room.ID, uid).Count(&rows).Error
	if err != nil {
		return -1, err
	}
	return int(rows), nil
}

func (repo *RoomRepository) Create(db *gorm.DB, room *models.Room) error {
	err := db.Create(room).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *RoomRepository) GetByID(db *gorm.DB, uid string) ([]models.Room, error) {
	var result []models.Room
	err := db.Table("user_rooms").Where("user_rooms.user_id = ?", uid).Joins("join rooms on rooms.id = user_rooms.room_id").Where("is_private = ?", false).Scan(&result).Error
	if err != nil {
		return []models.Room{}, err
	}
	return result, nil
}

func (repo *RoomRepository) Peek(db *gorm.DB, room *models.Room) error {
	var result models.Room
	return db.Table("messages").Where("room_id = ?", room.ID).Find(&(result.Messages)).Error
}
