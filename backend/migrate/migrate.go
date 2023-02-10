package main

import (
	"github.com/alkuinvito/malakh-api/initializers"
	"github.com/alkuinvito/malakh-api/models"
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
