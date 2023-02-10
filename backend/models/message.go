package models

import "time"

type Message struct {
	ID uint
	Body   string
	UserID uint
	RoomID uint
	CreatedAt time.Time
}
