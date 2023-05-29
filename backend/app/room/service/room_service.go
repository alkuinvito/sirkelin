package service

import (
	"sirkelin/backend/app/room/repository"
	"sirkelin/backend/models"
)

type RoomService struct {
	repository repository.RoomRepository
}

type IRoomService interface {
	CreateRoom([]*models.User) error
}

func NewRoomService(repository repository.RoomRepository) *RoomService {
	return &RoomService{
		repository: repository,
	}
}

func (service *RoomService) CreateRoom([]*models.User) error {
	return nil
}
