package models

import (
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
