package models

import (
	"github.com/alkuinvito/sirkelin/initializers"
	"github.com/google/uuid"
	"time"
)

type Message struct {
	ID        string
	Body      string
	UserID    string
	RoomID    string
	CreatedAt time.Time
}

type SendMessageParams struct {
	Body string `json:"body"`
}

func InsertMessage(message *Message) error {
	message.ID = uuid.New().String()
	return initializers.DB.Create(&message).Error
}
