package service

import (
	"sirkelin/backend/app/room/repository"
	"sirkelin/backend/initializers"
	"sirkelin/backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomService struct {
	repository repository.RoomRepository
	db *gorm.DB
}

type IRoomService interface {
	Create([]*models.User) (string, error)
}

func NewRoomService(repository repository.RoomRepository) *RoomService {
	return &RoomService{
		repository: repository,
	}
}

func (service *RoomService) Create([]*models.User) (string, error) {
	roomID := uuid.NewString()

	tx := service.db.Begin()
	defer initializers.CommitOrRollback(tx)
	err := service.repository.Create(tx, &models.Room{ID: roomID})
	if err != nil {
		return "", err
	}

	return roomID, nil
}
