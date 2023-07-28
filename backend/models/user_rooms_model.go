package models

type UserRooms struct {
	RoomID string `gorm:"primaryKey;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID string `gorm:"primaryKey;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
