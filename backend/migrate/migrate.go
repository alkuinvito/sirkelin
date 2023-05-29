package main

import (
	"sirkelin/backend/initializers"
)

func init() {
	initializers.LoadEnvVar()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate((&models.User{}))
	initializers.DB.AutoMigrate((&models.Message{}))
	initializers.DB.AutoMigrate((&models.Room{}))
}
