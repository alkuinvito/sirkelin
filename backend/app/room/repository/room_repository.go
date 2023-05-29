package repository

import (
	"gorm.io/gorm"
	"sirkelin/backend/initializers"
	"sirkelin/backend/models"
)

type RoomRepository struct {
	db *gorm.DB
}

type IRoomRepository interface {
	Count(room *models.Room, uid string) (int, error)
	FindByID(uid string) ([]models.Room, error)
	Insert(room *models.Room) error
	Peek(room *models.Room) error
}

func NewRoomRepository() *RoomRepository {
	return &RoomRepository{db: initializers.DB}
}

func (repo *RoomRepository) Count(room *models.Room, uid string) (int, error) {
	var rows int64
	err := repo.db.Table("user_rooms").Where("room_id = ? AND user_id = ?", room.ID, uid).Count(&rows).Error
	if err != nil {
		return -1, err
	}
	return int(rows), nil
}

func (repo *RoomRepository) FindByID(uid string) ([]models.Room, error) {
	var result []models.Room
	err := repo.db.Table("user_rooms").Where("user_rooms.user_id = ?", uid).Joins("join rooms on rooms.id = user_rooms.room_id").Where("is_private = ?", false).Scan(&result).Error
	if err != nil {
		return []models.Room{}, err
	}
	return result, nil
}

func (repo *RoomRepository) Insert(room *models.Room) error {
	err := repo.db.Create(room).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *RoomRepository) Peek(room *models.Room) error {
	var result models.Room
	return repo.db.Table("messages").Where("room_id = ?", room.ID).Find(&(result.Messages)).Error
}
