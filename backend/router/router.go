package router

import (
	"os"
	roomController "sirkelin/backend/app/room/controller"
	userController "sirkelin/backend/app/user/controller"
	"sirkelin/backend/app/websocket"
	"sirkelin/backend/middlewares"

	"github.com/gin-gonic/gin"
)

type Router struct {
	middleware *middlewares.Middleware
	user       *userController.UserController
	room       *roomController.RoomController
	hub        *websocket.Hub
}

func NewRouter(middleware *middlewares.Middleware, auth *userController.UserController, room *roomController.RoomController, hub *websocket.Hub) *Router {
	return &Router{
		middleware: middleware,
		user:       auth,
		room:       room,
		hub:        hub,
	}
}

func (router *Router) Handle() *gin.Engine {
	handler := gin.Default()
	gin.SetMode(os.Getenv("APP_MODE"))

	go router.hub.Run()

	authGroup := handler.Group("/auth")
	{
		authGroup.POST("/sign-in", router.user.SignIn)
		authGroup.POST("/sign-out", router.user.SignOut)
	}

	roomGroup := handler.Group("/room")
	{
		roomGroup.Use(router.middleware.AuthenticatedUser())
		roomGroup.POST("/", router.room.CreateRoom)
		roomGroup.GET("/private", router.room.GetPrivateRooms)
		roomGroup.GET("/", router.room.GetRooms)

		roomGroup.Use(router.middleware.AuthorizedUser())
		roomGroup.DELETE("/:id", router.room.Delete)
		roomGroup.GET("/:id", router.room.GetRoomById)
		roomGroup.POST("/:id/message", router.room.PushMessage)
		roomGroup.PUT("/:id", router.room.UpdateRoom)
	}

	userGroup := handler.Group("/user")
	{
		userGroup.Use(router.middleware.AuthenticatedUser())
		userGroup.GET("/", router.user.GetAll)
		userGroup.GET("/:id", router.user.GetByID)
		userGroup.PUT("/:id", router.user.UpdateProfile)
	}

	return handler
}
