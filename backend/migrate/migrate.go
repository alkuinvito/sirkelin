package main

import (
	"github.com/alkuinvito/sirkelin/initializers"
	"github.com/alkuinvito/sirkelin/models"
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
