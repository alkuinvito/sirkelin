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
	{
		authGroup.POST("/sign-in", controllers.SignIn)
	}

	privateGroup := router.Group("/private")
	{
		privateGroup.Use(middlewares.RoomAccess())
		privateGroup.GET("/", controllers.GetPrivateList)
		privateGroup.POST("create", controllers.CreatePrivateRoom)
	}

	roomGroup := router.Group("/room")
	{
		roomGroup.Use(middlewares.RoomAccess())
		roomGroup.GET("/", controllers.GetRoomList)
		roomGroup.POST("/create", controllers.CreateRoom)
		messageHandler := roomGroup.Group("/:id")
		{
			messageHandler.Use(middlewares.RoomPrivillege())
			messageHandler.POST("/", controllers.SendMessage)
			messageHandler.GET("/", controllers.GetMessages)
		}
	}

	return router
}