package repository

import (
	"github.com/alkuinvito/sirkelin/app/auth/repository"
	"gorm.io/gorm"
	"time"

	"github.com/alkuinvito/sirkelin/initializers"
)

type Room struct {
	ID        string
	Name      string
	Picture   string
	Users     []*repository.User `gorm:"many2many:user_rooms"`
	Messages  []Message
	IsPrivate bool `gorm:"not null"`
	CreatedAt time.Time
}

type CreateRoomParams struct {
	Name      string             `json:"name"`
	Picture   string             `json:"picture"`
	Users     []*repository.User `json:"users"`
	IsPrivate bool               `json:"is_private"`
}

type RoomRepository struct {
	db *gorm.DB
}

type IRoomRepository interface {
	Count(room *Room, uid string) (int, error)
	FindByID(uid string) ([]Room, error)
	Insert(room *Room) error
	Peek(room *Room) error
}

func Init() *RoomRepository {
	return &RoomRepository{db: initializers.DB}
}

func (repo *RoomRepository) Count(room *Room, uid string) (int, error) {
	var rows int64
	err := repo.db.Table("user_rooms").Where("room_id = ? AND user_id = ?", room.ID, uid).Count(&rows).Error
	if err != nil {
		return -1, err
	}
	return int(rows), nil
}

func (repo *RoomRepository) FindByID(uid string) ([]Room, error) {
	var result []Room
	err := repo.db.Table("user_rooms").Where("user_rooms.user_id = ?", uid).Joins("join rooms on rooms.id = user_rooms.room_id").Where("is_private = ?", false).Scan(&result).Error
	if err != nil {
		return []Room{}, err
	}
	return result, nil
}

func (repo *RoomRepository) Insert(room *Room) error {
	err := repo.db.Create(room).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *RoomRepository) Peek(room *Room) error {
	var result Room
	return repo.db.Table("messages").Where("room_id = ?", room.ID).Find(&(result.Messages)).Error
}
