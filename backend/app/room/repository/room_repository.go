package repository

import (
	"log"
	"sirkelin/backend/models"

	"gorm.io/gorm"
)

type RoomRepository struct {
}

type IRoomRepository interface {
	Count(db *gorm.DB, roomID, uid string) (int, error)
	Create(db *gorm.DB, room *models.Room, users []*models.User) error
	Delete(db *gorm.DB, id string) error
	GetPrivateRoomByParticipants(db *gorm.DB, ids []string) (models.Room, error)
	GetPrivateRooms(db *gorm.DB, uid string) ([]models.RoomList, error)
	GetRoomById(db *gorm.DB, id string) (models.Room, error)
	GetRooms(db *gorm.DB, uid string) ([]models.RoomList, error)
	Peek(db *gorm.DB, room *models.Room) error
	PushMessage(db *gorm.DB, message *models.Message) error
	Update(db *gorm.DB, room *models.Room) error
}

func NewRoomRepository() *RoomRepository {
	return &RoomRepository{}
}

func (repo *RoomRepository) Count(db *gorm.DB, roomID, uid string) (int, error) {
	var rows int64
	err := db.Table("user_rooms").Where("room_id = ? AND user_id = ?", roomID, uid).Count(&rows).Error
	if err != nil {
		return -1, err
	}
	return int(rows), nil
}

func (repo *RoomRepository) Create(db *gorm.DB, room *models.Room, users []*models.User) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var err error

		err = tx.Create(room).Error
		if err != nil {
			return err
		}

		for _, user := range users {
			err = tx.Model(user).Association("Rooms").Append(room)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (repo *RoomRepository) Delete(db *gorm.DB, id string) error {
	return db.Delete(&models.Room{ID: id}).Error
}

func (repo *RoomRepository) GetPrivateRoomByParticipants(db *gorm.DB, ids []string) (models.Room, error) {
	var result models.Room

	otherRoom := db.Table("user_rooms").Select("room_id").Where("user_id = ?", ids[1])
	err := db.Table("user_rooms").Select("user_rooms.room_id as id").Joins("join rooms on rooms.id = user_rooms.room_id").Where("user_id = ? and is_private = ? and room_id in (?)", ids[0], true, otherRoom).Find(&result).Error
	log.Println(err, result)
	if err != nil {
		return models.Room{}, err
	}
	return result, nil
}

func (repo *RoomRepository) GetPrivateRooms(db *gorm.DB, uid string) ([]models.RoomList, error) {
	var result []models.RoomList
	err := db.Raw("SELECT room_id as id, users.fullname as name, users.picture, is_private FROM user_rooms JOIN rooms ON rooms.id = user_rooms.room_id JOIN users ON users.id = user_rooms.user_id WHERE user_id <> ? AND rooms.is_private = TRUE AND room_id IN (SELECT room_id FROM user_rooms WHERE user_id = ?)", uid, uid).Scan(&result).Error
	if err != nil {
		return []models.RoomList{}, err
	}
	return result, nil
}

func (repo *RoomRepository) GetRoomById(db *gorm.DB, id string) (models.Room, error) {
	var result models.Room
	err := db.Where("id = ?", id).Preload("Users", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Fullname", "Picture")
	}).Preload("Messages").First(&result).Error
	if err != nil {
		return models.Room{}, err
	}
	return result, nil
}

func (repo *RoomRepository) GetRooms(db *gorm.DB, uid string) ([]models.RoomList, error) {
	var result []models.RoomList
	err := db.Table("user_rooms").Where("user_rooms.user_id = ?", uid).Joins("join rooms on rooms.id = user_rooms.room_id").Where("is_private = ?", false).Scan(&result).Error
	if err != nil {
		return []models.RoomList{}, err
	}
	return result, nil
}

func (repo *RoomRepository) Peek(db *gorm.DB, room *models.Room) error {
	var result models.Room
	return db.Table("messages").Where("room_id = ?", room.ID).Find(&(result.Messages)).Error
}

func (repo *RoomRepository) PushMessage(db *gorm.DB, message *models.Message) error {
	return db.Create(message).Error
}

func (repo *RoomRepository) Update(db *gorm.DB, room *models.Room) error {
	return db.Save(&room).Error
}
