package models

import "github.com/alkuinvito/sirkelin/initializers"

type PrivateResult struct {
	RoomId   uint
	UserId   uint
	Username string
	Picture  string
}

func PrivateList(userId string) []PrivateResult {
	var privateList []PrivateResult

	err := initializers.DB.Raw("SELECT room_id, user_id, users.username, users.picture FROM user_rooms JOIN rooms ON rooms.id = user_rooms.room_id JOIN users ON users.id = user_rooms.user_id WHERE user_id <> ? AND rooms.is_private = TRUE AND room_id IN (SELECT room_id FROM user_rooms WHERE user_id = ?)", userId, userId).Scan(&privateList)
	if err != nil {
		return privateList
	}

	return privateList
}
