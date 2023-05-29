package service

import (
	authRepository "github.com/alkuinvito/sirkelin/app/auth/repository"
	roomRepository "github.com/alkuinvito/sirkelin/app/room/repository"
)

type service struct {
	repository *roomRepository.RoomRepository
}

type IService interface {
	CreateRoom([]*authRepository.User, ...any)
}

func (svc *service) CreateRoom() {

}
