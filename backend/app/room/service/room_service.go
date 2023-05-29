package service

import (
	authRepository "sirkelin/backend/app/auth/repository"
	roomRepository "sirkelin/backend/app/room/repository"
)

type RoomService struct {
	repository roomRepository.RoomRepository
}

type IRoomService interface {
	CreateRoom([]*authRepository.User, ...any)
}

func NewRoomService(repository roomRepository.RoomRepository) *RoomService {
	return &RoomService{
		repository: repository,
	}
}

func (service *RoomService) CreateRoom() {

}
