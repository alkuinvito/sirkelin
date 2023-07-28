package main

import (
	"sirkelin/backend/initializers"
	"sirkelin/backend/models"
)

func init() {
	initializers.LoadEnvVar()
}

func main() {
	db := initializers.NewDB()
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Message{})
	db.AutoMigrate(&models.Room{})
	db.AutoMigrate(&models.UserRooms{})
}
