package router

import (
	"github.com/alkuinvito/sirkelin/controllers"
	"github.com/alkuinvito/sirkelin/middlewares"
	"github.com/gin-gonic/gin"
	"os"
)

func Handle() *gin.Engine {
	router := gin.Default()
	gin.SetMode(os.Getenv("APP_MODE"))

	authGroup := router.Group("/auth")
	controllers.AuthHandler(authGroup)

	privateGroup := router.Group("/private")
	privateGroup.Use(middlewares.RoomAccess())
	controllers.PrivateHandler(privateGroup)

	roomGroup := router.Group("/room")
	roomGroup.Use(middlewares.RoomAccess())
	controllers.RoomHandler(roomGroup)

	return router
}
