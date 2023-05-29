package main

import (
	"sirkelin/backend/initializers"
	"sirkelin/backend/models"
)

func init() {
	initializers.LoadEnvVar()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.User{})
	initializers.DB.AutoMigrate(&models.Message{})
	initializers.DB.AutoMigrate(&models.Room{})
}
