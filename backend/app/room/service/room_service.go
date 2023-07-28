package service

import (
	"log"
	"sirkelin/backend/app/room/repository"
	"sirkelin/backend/initializers"
	"sirkelin/backend/models"
	"sirkelin/backend/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomService struct {
	repository *repository.RoomRepository
	db         *gorm.DB
}

type IRoomService interface {
	CheckRoomParticipant(roomID, uid string) (bool, error)
	Create(req *models.CreateRoomParams) (string, error)
	Delete(id string) error
	GetPrivateRooms(uid string) ([]models.RoomList, error)
	GetRoomById(id string) (models.Room, error)
	GetRooms(uid string) ([]models.RoomList, error)
	PushMessage(uid, roomId string, params models.SendMessageParams) (string, error)
	UpdateRoom(id string, data models.UpdateRoomSchema) error
}

func NewRoomService(repository *repository.RoomRepository, db *gorm.DB) *RoomService {
	return &RoomService{
		repository: repository,
		db:         db,
	}
}

func (service *RoomService) CheckRoomParticipant(roomID, uid string) (bool, error) {
	tx := service.db
	rows, err := service.repository.Count(tx, roomID, uid)
	if err != nil {
		return false, err
	}

	return rows == 1, nil
}

func (service *RoomService) Create(req *models.CreateRoomParams) (string, error) {
	var err error

	roomID := uuid.NewString()
	room := &models.Room{ID: roomID, IsPrivate: req.IsPrivate}

	if req.IsPrivate {
		if len(req.Users) != 2 {
			return "", utils.ErrPrivateParticipantsNumber
		}

		tx := service.db
		var ids []string
		for _, user := range req.Users {
			ids = append(ids, user.ID)
		}
		result, err := service.repository.GetPrivateRoomByParticipants(tx, ids)
		if err != nil {
			return "", err
		}

		if result.ID != "" {
			return result.ID, nil
		}
	} else {
		if len(req.Users) < 2 {
			return "", utils.ErrMinimumParticipant
		}

		room.Name = req.Name
		room.Picture = req.Picture
	}

	tx := service.db.Begin()
	defer initializers.CommitOrRollback(tx)
	err = service.repository.Create(tx, room, req.Users)
	log.Println(err)
	if err != nil {
		return "", err
	}
	return roomID, nil
}

func (service *RoomService) Delete(id string) error {
	tx := service.db.Begin()
	defer initializers.CommitOrRollback(tx)
	return service.repository.Delete(tx, id)
}

func (service *RoomService) GetPrivateRooms(uid string) ([]models.RoomList, error) {
	tx := service.db
	rooms, err := service.repository.GetPrivateRooms(tx, uid)
	if err != nil {
		return []models.RoomList{}, err
	}
	return rooms, nil
}

func (service *RoomService) GetRoomById(id string) (models.Room, error) {
	tx := service.db
	room, err := service.repository.GetRoomById(tx, id)
	if err != nil {
		return models.Room{}, err
	}
	return room, nil
}

func (service *RoomService) GetRooms(uid string) ([]models.RoomList, error) {
	tx := service.db
	rooms, err := service.repository.GetRooms(tx, uid)
	if err != nil {
		return []models.RoomList{}, err
	}
	return rooms, nil
}

func (service *RoomService) PushMessage(uid, roomId string, params models.SendMessageParams) (string, error) {
	messageID := uuid.NewString()
	message := &models.Message{
		ID:        messageID,
		Body:      params.Body,
		UserID:    uid,
		RoomID:    roomId,
		CreatedAt: time.Time{},
	}

	tx := service.db.Begin()
	defer initializers.CommitOrRollback(tx)
	err := service.repository.PushMessage(tx, message)
	if err != nil {
		return "", err
	}

	return messageID, nil
}

func (service *RoomService) UpdateRoom(id string, data models.UpdateRoomSchema) error {
	room := models.Room{
		ID:      id,
		Name:    data.Name,
		Picture: data.Picture,
	}
	tx := service.db.Begin()
	defer initializers.CommitOrRollback(tx)
	err := service.repository.Update(tx, &room)
	if err != nil {
		return err
	}
	return nil
}
