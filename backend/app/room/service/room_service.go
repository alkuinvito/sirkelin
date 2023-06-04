package service

import (
	"sirkelin/backend/app/room/repository"
	"sirkelin/backend/initializers"
	"sirkelin/backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomService struct {
	repository *repository.RoomRepository
	db         *gorm.DB
}

type IRoomService interface {
	Create(users []*models.User) (string, error)
	GetByUID(uid string) ([]models.Room, error)
}

func NewRoomService(repository *repository.RoomRepository, db *gorm.DB) *RoomService {
	return &RoomService{
		repository: repository,
		db:         db,
	}
}

func (service *RoomService) Create(users []*models.User) (string, error) {
	roomID := uuid.NewString()
	room := &models.Room{
		ID:    roomID,
		Users: users,
	}

	tx := service.db.Begin()
	defer initializers.CommitOrRollback(tx)
	err := service.repository.Create(tx, room)
	if err != nil {
		return "", err
	}
	return roomID, nil
}

func (service *RoomService) GetByUID(uid string) ([]models.Room, error) {
	tx := service.db.Begin()
	defer initializers.CommitOrRollback(tx)
	rooms, err := service.repository.GetByUID(tx, uid)
	if err != nil {
		return []models.Room{}, err
	}
	return rooms, nil
}
