package models

import "time"

type Message struct {
	ID        uint
	Body      string
	UserID    string
	RoomID    uint
	CreatedAt time.Time
}
