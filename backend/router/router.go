package router

import (
	authController "github.com/alkuinvito/sirkelin/app/auth/controller"
	"os"

	"github.com/alkuinvito/sirkelin/controllers"
	"github.com/alkuinvito/sirkelin/middlewares"
	"github.com/gin-gonic/gin"
)

func Handle() *gin.Engine {
	router := gin.Default()
	gin.SetMode(os.Getenv("APP_MODE"))

	ac := authController.Init()

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/sign-in", ac.SignIn)
		authGroup.POST("/sign-out", ac.SignOut)
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
			messageHandler.Use(middlewares.RoomPrivilege())
			messageHandler.POST("/", controllers.SendMessage)
			messageHandler.GET("/", controllers.GetMessages)
		}
	}

	userGroup := router.Group("/user")
	{
		userGroup.Use(middlewares.RoomAccess())
		userGroup.GET("/list", controllers.GetUsers)
	}

	return router
}
